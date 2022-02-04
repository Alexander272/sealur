package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initProRoutes(api *gin.RouterGroup) {
	pro := api.Group("/sealur-pro")
	{
		pro.GET("/ping", h.notImplemented)

		stand := pro.Group("/standards")
		{
			stand.GET("/", h.notImplemented)
			stand.POST("/", h.notImplemented)
			stand.PUT("/:id", h.notImplemented)
			stand.DELETE("/:id", h.notImplemented)
		}
	}
}
