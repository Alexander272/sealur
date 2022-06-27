package middleware

import (
	"fmt"
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

	isRefresh, err := m.services.Session.CheckSession(c, user, token)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}

	if isRefresh {
		token, err := m.services.Session.SignIn(c, user)
		if err != nil {
			models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "failed to refresh session")
			return
		}

		c.SetCookie(m.CookieName, token, int(m.auth.RefreshTokenTTL.Seconds()), "/", m.auth.Domain, m.auth.Secure, true)
	}

	c.Set(userIdCtx, user.Id)
	for _, r := range user.Roles {
		c.Set(fmt.Sprintf("%s_%s", userRolesCtx, r.Service), r.Role)
	}
	// c.Set(userRolesCtx, user.Roles)
}

func (m *Middleware) AccessForSuperUser(c *gin.Context) {
	userId, _ := c.Get(userIdCtx)
	role, exists := c.Get(fmt.Sprintf("%s_sealur", userRolesCtx))
	if !exists {
		models.NewErrorResponse(c, http.StatusUnauthorized, "roles empty", "failed to get role")
	}

	if role != "superuser" {
		models.NewErrorResponse(c, http.StatusForbidden, role.(string), "access not allowed")
		logger.Error(userId)
	}
}
