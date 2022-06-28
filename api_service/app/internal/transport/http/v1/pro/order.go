package pro

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	order := api.Group("/orders", h.middleware.UserIdentity)
	{
		order.GET("/all", h.getAllOrders)
		order.POST("/", h.createOrder)
		order.POST("/:orderId", h.saveOrder)
		order.DELETE("/:orderId", h.deleteOrder)

		//TODO добавить ендпоинт для копирования заказа или его позиций в новый заказ

		order.GET("/positions", h.getCurPosition)
		positions := order.Group("/:orderId/positions")
		{
			positions.GET("/", h.getPosition)
			positions.POST("/", h.addPosition)
			positions.PATCH("/:id", h.updatePosition)
			positions.DELETE("/:id", h.removePosition)
		}
	}
}

// @Summary Get All Orders
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description получение всех заказов
// @ModuleID getAllOrders
// @Accept json
// @Produce json
// @Param userId query string true "user id"
// @Success 200 {object} models.DataResponse{Data=[]proto.Order}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/all [get]
func (h *Handler) getAllOrders(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty userId param")
		return
	}

	order, err := h.proClient.GetAllOrders(c, &proto.GetAllOrdersRequest{
		UserId: userId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: order.Orders, Count: len(order.Orders)})
}

// @Summary Create Order
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description создание заказа
// @ModuleID createOrder
// @Accept json
// @Produce json
// @Param order body models.OrderDTO true "order info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/ [post]
func (h *Handler) createOrder(c *gin.Context) {
	var dto models.OrderDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	if dto.Id == "" {
		id, err := uuid.NewUUID()
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "failed to generate group id")
			return
		}
		dto.Id = id.String()
	}

	order, err := h.proClient.CreateOrder(c, &proto.CreateOrderRequest{
		OrderId: dto.Id,
		Count:   dto.Count,
		UserId:  dto.UserId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/orders/%s", order.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Id: order.Id, Message: "Created"})
}

// @Summary Save Order
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description сохранение заказа
// @ModuleID saveOrder
// @Accept json
// @Produce json
// @Param orderId path string true "order id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/{orderId} [post]
func (h *Handler) saveOrder(c *gin.Context) {

	c.JSON(http.StatusOK, models.IdResponse{Message: "Saved"})
}

// @Summary Delete Order
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description сохранение заказа
// @ModuleID deleteOrder
// @Accept json
// @Produce json
// @Param orderId path string true "order id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/{orderId} [delete]
func (h *Handler) deleteOrder(c *gin.Context) {
	id := c.Param("orderId")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty id", "empty id param")
		return
	}

	order, err := h.proClient.DeleteOrder(c, &proto.DeleteOrderRequest{OrderId: id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: order.Id, Message: "Deleted"})
}

// @Summary Get Position
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description получение позиций к заказу
// @ModuleID getPosition
// @Accept json
// @Produce json
// @Param orderId path string true "order id"
// @Success 200 {object} models.DataResponse{Data=[]proto.OrderPosition}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/{orderId}/positions [get]
func (h *Handler) getPosition(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty orderId param")
		return
	}

	positions, err := h.proClient.GetPositions(c, &proto.GetPositionsRequest{
		OrderId: orderId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: positions.Positions, Count: len(positions.Positions)})
}

// @Summary Get Cur Position
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description получение позиций к несохранненому заказу
// @ModuleID getCurPosition
// @Accept json
// @Produce json
// @Param userId query string true "user id"
// @Success 200 {object} models.DataResponse{Data=[]proto.OrderPosition}
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/positions [get]
func (h *Handler) getCurPosition(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty userId param")
		return
	}

	positions, err := h.proClient.GetCurPositions(c, &proto.GetCurPositionsRequest{
		UserId: userId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: positions.Positions, Count: len(positions.Positions)})
}

// @Summary Add Position
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description добавление позиции к заказу
// @ModuleID addPosition
// @Accept json
// @Produce json
// @Param orderId path string true "order id"
// @Param position body models.PositionDTO true "position info"
// @Success 201 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/{orderId}/positions [post]
func (h *Handler) addPosition(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty orderId param")
		return
	}
	var dto models.PositionDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	count, err := strconv.Atoi(dto.Count)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "", "invalid data send")
		return
	}

	positions, err := h.proClient.AddPosition(c, &proto.AddPositionRequest{
		OrderId:     orderId,
		Designation: dto.Designation,
		Description: dto.Descriprion,
		Count:       int32(count),
		Sizes:       dto.Sizes,
		Drawing:     dto.Drawing,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: positions.Id, Message: "Created"})
}

// @Summary Update Position
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description обновление позиции к заказу
// @ModuleID updatePosition
// @Accept json
// @Produce json
// @Param orderId path string true "order id"
// @Param id path string true "position id"
// @Param position body models.PositionDTO true "position info"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/{orderId}/positions/{id} [patch]
func (h *Handler) updatePosition(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty orderId param")
		return
	}
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	var dto models.PositionDTO
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	count, err := strconv.Atoi(dto.Count)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "", "invalid data send")
		return
	}

	positions, err := h.proClient.UpdatePosition(c, &proto.UpdatePositionRequest{
		Id:    id,
		Count: int32(count),
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: positions.Id, Message: "Updated"})
}

// @Summary Remove Position
// @Tags Sealur Pro -> orders
// @Security ApiKeyAuth
// @Description удаление позиции к заказу
// @ModuleID removePosition
// @Accept json
// @Produce json
// @Param orderId path string true "order id"
// @Param id path string true "position id"
// @Success 200 {object} models.IdResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /sealur-pro/orders/{orderId}/positions/{id} [delete]
func (h *Handler) removePosition(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty orderId param")
		return
	}
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	positions, err := h.proClient.RemovePosition(c, &proto.RemovePositionRequest{
		OrderId: orderId,
		Id:      id,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Id: positions.Id, Message: "Deleted"})
}
