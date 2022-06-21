package user

import "github.com/gin-gonic/gin"

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.singUp)
		auth.POST("/sign-out", h.signOut)
	}

	users := api.Group("/users")
	{
		users.GET("/all", h.getAllUsers)
		users.GET("/new", h.getNewUsers)
		users.GET("/:id", h.getUser)
		users.PATCH("/:id", h.updateUser)
		users.DELETE("/:id", h.deleteUser)
		users.POST("/confirm", h.confirmUser)
		roles := users.Group("/roles")
		{
			roles.POST("/", h.createRole)
			roles.PUT("/:id", h.updateRole)
			roles.DELETE("/:id", h.deleteRole)
		}
	}
}
