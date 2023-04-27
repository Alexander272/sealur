package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/models/user_model"
	"github.com/Alexander272/sealur/api_service/internal/service"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userApi    user_api.UserServiceClient
	emailApi   email_api.EmailServiceClient
	auth       config.AuthConfig
	http       config.HttpConfig
	services   *service.Services
	cookieName string
}

func NewAuthHandler(
	userApi user_api.UserServiceClient, emailApi email_api.EmailServiceClient,
	auth config.AuthConfig, http config.HttpConfig,
	services *service.Services,
	cookieName string,
) *AuthHandler {
	return &AuthHandler{
		userApi:    userApi,
		emailApi:   emailApi,
		auth:       auth,
		http:       http,
		services:   services,
		cookieName: cookieName,
	}
}

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	handler := NewAuthHandler(h.userApi, h.emailApi, h.auth, h.http, h.services, h.cookieName)

	auth := api.Group("/auth")
	{
		auth.POST("/sign-in", handler.signIn)
		auth.POST("/sign-up", handler.singUp)
		auth.POST("/sign-out", handler.signOut)
		auth.POST("/refresh", handler.refresh)
	}
}

func (h *AuthHandler) signIn(c *gin.Context) {
	var req user_model.SignIn
	if err := c.BindJSON(&req); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Введены некорректные данные")
		return
	}

	dto := &user_api.GetUserByEmail{
		Email:    req.Email,
		Password: req.Password,
	}

	limit, err := h.services.Limit.Get(c, c.ClientIP())
	if err != nil && !errors.Is(err, models.ErrClientIPNotFound) {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error()+" user email: "+dto.Email, "Произошла ошибка")
		return
	}
	if errors.Is(err, models.ErrClientIPNotFound) {
		h.services.Limit.Create(c, c.ClientIP())
	}

	//TODO
	// if limit.Count == h.auth.CountAttempt {
	// 	h.emailApi.SendBlocked(c, &email_api.BlockedUserRequest{Ip: c.ClientIP(), Login: dto.Login})
	// }

	if limit.Count >= h.auth.CountAttempt {
		h.services.Limit.AddAttempt(c, c.ClientIP())
		models.NewErrorResponse(
			c, http.StatusTooManyRequests,
			fmt.Sprintf("too many request (%d >= %d)", limit.Count, h.auth.CountAttempt),
			fmt.Sprintf("Много некорректных запросов. Доступ заблокирован на %.0f минут", h.auth.LimitAuthTTL.Minutes()),
		)
		return
	}

	user, err := h.userApi.GetByEmail(c, dto)
	if err != nil && !strings.Contains(err.Error(), "invalid credentials") {
		if strings.Contains(err.Error(), "user not verified") {
			models.NewErrorResponse(c, http.StatusBadRequest, "user not verified",
				"Учетная запись не активирована. Для активации учетной записи перейдите по ссылке, отправленной вам в письме.",
			)
			return
		}
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")
		return
	}

	if user == nil {
		h.services.Limit.AddAttempt(c, c.ClientIP())
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid credentials", "Введены некорректные данные")
		return
	}
	h.services.Limit.Remove(c, c.ClientIP())

	// запись в редисе и генерация токенов
	token, err := h.services.Session.SignIn(c, user)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")
		return
	}

	// _, err = h.userClient.AddIp(c, &user_api.AddIpRequest{
	// 	UserId: user.User.Id,
	// 	Ip:     c.ClientIP(),
	// })
	// if err != nil {
	// 	logger.Error(err)
	// }

	c.SetCookie(h.cookieName, token, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.Host, h.auth.Secure, true)
	c.JSON(http.StatusOK, models.DataResponse{Data: user})
}

func (h *AuthHandler) singUp(c *gin.Context) {
	var dto *user_api.CreateUser
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Введены некорректные данные")
		return
	}

	id, err := h.userApi.Create(c, dto)
	if err != nil {
		// if errors.Is(err, models.ErrUserExist) {
		if strings.Contains(err.Error(), models.ErrUserExist.Error()) {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Пользователь с таким email уже зарегистрирован")
			return
		}
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")
		return
	}

	// генерировать код для подтверждения и записывать его в редис (с id пользователя)
	code, err := h.services.Confirm.Create(c, id.Id)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка")
		return
	}

	logger.Debug(fmt.Sprintf("%s/auth/confirm?code=%s", h.http.Link, code))

	data := &email_api.ConfirmUserRequest{
		Name:  dto.Name,
		Email: dto.Email,
		Link:  fmt.Sprintf("%s/auth/confirm?code=%s", h.http.Link, code),
	}
	_, err = h.emailApi.ConfirmUser(c, data)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка при отправлении письма")
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Registration completed successfully"})
}

func (h *AuthHandler) signOut(c *gin.Context) {
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

	c.SetCookie(h.cookieName, "", -1, "/", c.Request.Host, h.auth.Secure, true)
	c.JSON(http.StatusNoContent, models.IdResponse{Message: "Sign-out completed successfully"})
}

func (h *AuthHandler) refresh(c *gin.Context) {
	token, err := c.Cookie(h.cookieName)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}
	if token == "" {
		models.NewErrorResponse(c, http.StatusUnauthorized, "empty token", "user is not authorized")
		return
	}

	user, err := h.services.Session.TokenParse(token)
	if err != nil {
		c.SetCookie(h.cookieName, token, -1, "/", c.Request.Host, h.auth.Secure, true)
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error()+" token: "+token, "user is not authorized")
		return
	}

	_, err = h.services.CheckSession(c, user, token)
	if err != nil {
		c.SetCookie(h.cookieName, token, -1, "/", c.Request.Host, h.auth.Secure, true)
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error()+" token: "+token+" userId: "+user.Id, "user is not authorized")
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
