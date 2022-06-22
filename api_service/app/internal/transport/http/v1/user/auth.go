package user

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/proto_user"
	"github.com/gin-gonic/gin"
)

// @Summary SignIn
// @Tags auth
// @Description вход в систему
// @ModuleID signIn
// @Accept json
// @Produce json
// @Param data body models.SignIn true "credentials"
// @Success 200 {object} models.DataResponse{data=proto_user.User}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var dto models.SignIn
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	user, err := h.userClient.GetUser(c, &proto_user.GetUserRequest{Login: dto.Login, Password: dto.Password})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	if user.User == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid crenditails", "invalid data send")
		return
	}

	// запись в редисе и генерация токенов
	token, err := h.services.SignIn(c, user.User)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.SetCookie(h.cookieName, token, int(h.auth.RefreshTokenTTL.Seconds()), "/", h.auth.Domain, h.auth.Secure, true)
	c.JSON(http.StatusOK, models.DataResponse{Data: user.User})
}

// @Summary SignUp
// @Tags auth
// @Description регистрация
// @ModuleID singUp
// @Accept json
// @Produce json
// @Param data body models.SignUp true "user info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) singUp(c *gin.Context) {
	var dto models.SignUp
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	req := proto_user.CreateUserRequest{
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
// @Tags auth
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

	c.SetCookie(h.cookieName, token, 0, "/", h.auth.Domain, h.auth.Secure, true)
	c.JSON(http.StatusNoContent, models.IdResponse{Message: "Sign-out completed successfully"})
}

// @Summary Refresh
// @Tags auth
// @Description вход в систему (при обновлении страницы)
// @ModuleID refresh
// @Accept json
// @Produce json
// @Success 200 {object} models.DataResponse{data=proto_user.UserResponse}
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

	// запись в редисе и генерация токенов
	newToken, err := h.services.SignIn(c, user)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.SetCookie(h.cookieName, newToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", h.auth.Domain, h.auth.Secure, true)
	c.JSON(http.StatusOK, models.DataResponse{Data: user})
}
