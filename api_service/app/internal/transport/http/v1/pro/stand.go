package pro

import "github.com/gin-gonic/gin"

func (h *Handler) initStandRoutes(api *gin.RouterGroup) {
	stand := api.Group("/standards")
	{
		stand.GET("/", h.notImplemented)
		stand.POST("/", h.notImplemented)
		stand.PUT("/:id", h.notImplemented)
		stand.DELETE("/:id", h.notImplemented)
	}
}
