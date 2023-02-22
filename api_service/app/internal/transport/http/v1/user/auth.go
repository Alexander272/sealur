package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/user_model"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/user_api"
	"github.com/gin-gonic/gin"
)

// @Summary SignIn
// @Tags Auth
// @Description вход в систему
// @ModuleID signIn
// @Accept json
// @Produce json
// @Param data body user_model.SignIn true "credentials"
// @Success 200 {object} models.DataResponse{data=user_api.User}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var dto user_model.SignIn
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	limit, err := h.services.Limit.Get(c, c.ClientIP())
	if err != nil && !errors.Is(err, models.ErrClientIPNotFound) {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	if errors.Is(err, models.ErrClientIPNotFound) {
		h.services.Limit.Create(c, c.ClientIP())
	}

	if limit.Count == h.auth.CountAttempt {
		h.emailClient.SendBlocked(c, &email_api.BlockedUserRequest{Ip: c.ClientIP(), Login: dto.Login})
	}

	if limit.Count >= h.auth.CountAttempt {
		h.services.AddAttempt(c, c.ClientIP())
		models.NewErrorResponse(
			c, http.StatusTooManyRequests,
			fmt.Sprintf("too many request (%d >= %d)", limit.Count, h.auth.CountAttempt),
			"too many request",
		)
		return
	}

	user, err := h.userClient.GetUser(c, &user_api.GetUserRequest{Login: dto.Login, Password: dto.Password})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	if user.User == nil {
		h.services.AddAttempt(c, c.ClientIP())
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid crenditails", "invalid data send")
		return
	}
	h.services.Remove(c, c.ClientIP())

	// запись в редисе и генерация токенов
	token, err := h.services.SignIn(c, user.User)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	_, err = h.userClient.AddIp(c, &user_api.AddIpRequest{
		UserId: user.User.Id,
		Ip:     c.ClientIP(),
	})
	if err != nil {
		logger.Error(err)
	}

	c.SetCookie(h.cookieName, token, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.Host, h.auth.Secure, true)
	c.JSON(http.StatusOK, models.DataResponse{Data: user.User})
}

// @Summary SignUp
// @Tags Auth
// @Description регистрация
// @ModuleID singUp
// @Accept json
// @Produce json
// @Param data body user_model.SignUp true "user info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) singUp(c *gin.Context) {
	var dto user_model.SignUp
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	req := user_api.CreateUserRequest{
		Organization: dto.Organization,
		Name:         dto.Name,
		Email:        dto.Email,
		City:         dto.City,
		Position:     dto.Position,
		Phone:        dto.Phone,
	}

	_, err := h.userClient.CreateUser(c, &req)
	if err != nil {
		// if errors.Is(err, models.ErrFlangeAlreadyExists) {
		// 	models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
		// 	return
		// }
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Registration completed successfully"})
}

// @Summary SignOut
// @Tags Auth
// @Description выход из аккаунта
// @ModuleID signOut
// @Accept json
// @Produce json
// @Success 204 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/sign-out [post]
func (h *Handler) signOut(c *gin.Context) {
	token, err := c.Cookie(h.cookieName)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	user, err := h.services.Session.TokenParse(token)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	if err := h.services.Session.SingOut(c, user.Id); err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.SetCookie(h.cookieName, "", 0, "/", c.Request.Host, h.auth.Secure, true)
	c.JSON(http.StatusNoContent, models.IdResponse{Message: "Sign-out completed successfully"})
}

// @Summary Refresh
// @Tags Auth
// @Description вход в систему (при обновлении страницы)
// @ModuleID refresh
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=user_api.UserResponse}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	token, err := c.Cookie(h.cookieName)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	user, err := h.services.Session.TokenParse(token)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	_, err = h.services.CheckSession(c, user, token)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	// удаление записей из редиса
	if err := h.services.Session.SingOut(c, user.Id); err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	// запись в редисе и генерация токенов
	newToken, err := h.services.SignIn(c, user)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.SetCookie(h.cookieName, newToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.Host, h.auth.Secure, true)
	c.JSON(http.StatusOK, models.DataResponse{Data: user})
}
