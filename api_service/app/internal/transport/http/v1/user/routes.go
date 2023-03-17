package user

// func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
// auth := api.Group("/auth")
// {
// 	auth.POST("/sign-in", h.signIn)
// 	auth.POST("/sign-up", h.singUp)
// 	auth.POST("/sign-out", h.signOut)
// 	auth.POST("/refresh", h.refresh)
// }

// api.GET("users/all/:page", h.getAllUsers)
// users := api.Group("/users", h.middleware.UserIdentity)
// {
// 	// users.GET("/all/:page", h.middleware.AccessForSuperUser, h.getAllUsers)
// 	users.GET("/new", h.middleware.AccessForSuperUser, h.getNewUsers)
// 	users.GET("/:id", h.getUser)
// 	users.PATCH("/:id", h.updateUser)
// 	users.DELETE("/:id", h.deleteUser)
// 	users.POST("/confirm", h.middleware.AccessForSuperUser, h.confirmUser)
// 	users.POST("/clear", h.middleware.AccessForSuperUser, h.clearLimit)
// 	users.DELETE("/reject/:id", h.middleware.AccessForSuperUser, h.rejectUser)
// 	roles := users.Group("/roles", h.middleware.AccessForSuperUser)
// 	{
// 		roles.POST("/", h.createRole)
// 		roles.PUT("/:id", h.updateRole)
// 		roles.DELETE("/:id", h.deleteRole)
// 	}
// }
// }
