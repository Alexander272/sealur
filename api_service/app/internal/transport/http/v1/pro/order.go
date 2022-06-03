package pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	order := api.Group("/orders")
	{
		// order.GET("/", h.getPutg)
		// order.POST("/", h.createPutg)
		// order.PUT("/:id", h.updatePutg)
		// order.DELETE("/:id", h.deletePutg)

		order.POST("/", h.createOrder)
	}
}

// @Summary Create Order
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description создание заказа
// @ModuleID createOrder
// @Accept json
// @Produce json
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/ [post]
func (h *Handler) createOrder(c *gin.Context) {

	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}
