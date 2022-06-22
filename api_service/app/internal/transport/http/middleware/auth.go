package middleware

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	userIdCtx    = "userId"
	userRolesCtx = "roles"
)

func (m *Middleware) UserIdentity(c *gin.Context) {
	token, err := c.Cookie(m.CookieName)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	user, err := m.services.Session.TokenParse(token)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}
	logger.Debug("user: ", user)

	isRefresh, err := m.services.Session.CheckSession(c, user, token)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	logger.Debug("isRefresh ", isRefresh)
	if isRefresh {
		token, err := m.services.Session.SignIn(c, user)
		if err != nil {
			models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "failed to refresh session")
			return
		}

		c.SetCookie(m.CookieName, token, int(m.auth.RefreshTokenTTL.Seconds()), "/", m.auth.Domain, m.auth.Secure, true)
	}

	c.Set(userIdCtx, user.Id)
	c.Set(userRolesCtx, user.Roles)
}
