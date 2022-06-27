package user

import "github.com/gin-gonic/gin"

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.singUp)
		auth.POST("/sign-out", h.signOut)
		auth.POST("/refresh", h.refresh)
	}

	users := api.Group("/users", h.middleware.UserIdentity)
	{
		users.GET("/all", h.middleware.AccessForSuperUser, h.getAllUsers)
		users.GET("/new", h.middleware.AccessForSuperUser, h.getNewUsers)
		users.GET("/:id", h.getUser)
		users.PATCH("/:id", h.updateUser)
		users.DELETE("/:id", h.deleteUser)
		users.POST("/confirm", h.middleware.AccessForSuperUser, h.confirmUser)
		users.DELETE("/reject/:id", h.middleware.AccessForSuperUser, h.rejectUser)
		roles := users.Group("/roles", h.middleware.AccessForSuperUser)
		{
			roles.POST("/", h.createRole)
			roles.PUT("/:id", h.updateRole)
			roles.DELETE("/:id", h.deleteRole)
		}
	}
}
