package pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	order := api.Group("/order")
	{
		// order.GET("/", h.getPutg)
		// order.POST("/", h.createPutg)
		// order.PUT("/:id", h.updatePutg)
		// order.DELETE("/:id", h.deletePutg)

		order.POST("/drawing", h.createDrawing)
	}
}

// @Summary Create Drawing
// @Tags Sealur Pro -> order
// @Security ApiKeyAuth
// @Description создание чертежа
// @ModuleID createDrawing
// @Accept multipart/form-data
// @Produce json
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/order/drawing [post]
func (h *Handler) createDrawing(c *gin.Context) {

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}