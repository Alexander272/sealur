package middleware

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	UserIdCtx    = "userId"
	UserRolesCtx = "roles"
)

func (m *Middleware) UserIdentity(c *gin.Context) {
	token, err := c.Cookie(m.CookieName)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "user is not authorized")
		return
	}
	if token == "" {
		models.NewErrorResponse(c, http.StatusUnauthorized, "empty token", "user is not authorized")
		return
	}

	user, err := m.services.Session.TokenParse(token)
	if err != nil {
		//TODO проверить работоспособность
		c.SetCookie(m.CookieName, token, -1, "/", m.auth.Domain, m.auth.Secure, true)
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error()+" token: "+token, "user is not authorized")
		return
	}

	isRefresh, err := m.services.Session.CheckSession(c, user, token)
	if err != nil {
		c.SetCookie(m.CookieName, token, -1, "/", m.auth.Domain, m.auth.Secure, true)
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error()+" token: "+token+" userId: "+user.Id, "user is not authorized")
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

	c.Set(UserIdCtx, user.Id)
	// for _, r := range user.Roles {
	// 	c.Set(fmt.Sprintf("%s_%s", userRolesCtx, r.Service), r.Role)
	// }
	c.Set(UserRolesCtx, user.RoleCode)
}

func (m *Middleware) AccessForManager(c *gin.Context) {
	userId, _ := c.Get(UserIdCtx)
	role, exists := c.Get(UserRolesCtx)
	if !exists {
		models.NewErrorResponse(c, http.StatusUnauthorized, "roles empty", "failed to get role")
		return
	}

	if role != "manager" {
		models.NewErrorResponse(c, http.StatusForbidden, role.(string), "access not allowed")
		logger.Error(userId)
		return
	}
}

// TODO переписать Middleware
func (m *Middleware) AccessForProAdmin(c *gin.Context) {
	userId, _ := c.Get(UserIdCtx)
	role, exists := c.Get(fmt.Sprintf("%s_pro", UserRolesCtx))
	if !exists {
		models.NewErrorResponse(c, http.StatusUnauthorized, "roles empty", "failed to get role")
		return
	}

	if role != "admin" {
		models.NewErrorResponse(c, http.StatusForbidden, role.(string), "access not allowed")
		logger.Error(userId)
		return
	}
}

func (m *Middleware) AccessForMomentAdmin(c *gin.Context) {
	userId, _ := c.Get(UserIdCtx)
	role, exists := c.Get(fmt.Sprintf("%s_moment", UserRolesCtx))
	if !exists {
		models.NewErrorResponse(c, http.StatusUnauthorized, "roles empty", "failed to get role")
		return
	}

	if role != "admin" {
		models.NewErrorResponse(c, http.StatusForbidden, role.(string), "access not allowed")
		logger.Error(userId)
		return
	}
}

func (m *Middleware) AccessForSuperUser(c *gin.Context) {
	userId, _ := c.Get(UserIdCtx)
	role, exists := c.Get(fmt.Sprintf("%s_sealur", UserRolesCtx))
	if !exists {
		models.NewErrorResponse(c, http.StatusUnauthorized, "roles empty", "failed to get role")
		return
	}

	if role != "superuser" {
		models.NewErrorResponse(c, http.StatusForbidden, role.(string), "access not allowed")
		logger.Error(userId)
		return
	}
}
