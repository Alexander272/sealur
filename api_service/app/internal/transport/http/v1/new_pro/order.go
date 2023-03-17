package new_pro

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderApi order_api.OrderServiceClient
}

func NewOrderHandler(orderApi order_api.OrderServiceClient) *OrderHandler {
	return &OrderHandler{
		orderApi: orderApi,
	}
}

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	handler := NewOrderHandler(h.orderApi)

	order := api.Group("/orders", h.middleware.UserIdentity)
	{
		order.GET("/:id", handler.get)
		order.GET("/current", handler.getCurrent)
		order.GET("/all", handler.getAll)
		order.GET("/:id/заявка.zip", handler.getFile)
		order.POST("/", handler.create)
		order.POST("/save", handler.save)
	}
}

func (h *OrderHandler) get(c *gin.Context) {
	//TODO
}

func (h *OrderHandler) getCurrent(c *gin.Context) {
	userId, exists := c.Get(middleware.UserIdCtx)
	if !exists {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty user id param")
		return
	}

	order, err := h.orderApi.GetCurrent(c, &order_api.GetCurrentOrder{UserId: userId.(string)})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: order})
}

func (h *OrderHandler) getAll(c *gin.Context) {
	userId, exists := c.Get(middleware.UserIdCtx)
	if !exists {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty user id param")
		return
	}

	orders, err := h.orderApi.GetAll(c, &order_api.GetAllOrders{UserId: userId.(string)})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: orders.Orders, Count: len(orders.Orders)})
}

func (h *OrderHandler) getFile(c *gin.Context) {
	orderId := c.Param("id")
	if orderId == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty orderId param")
		return
	}

	stream, err := h.orderApi.GetFile(c, &order_api.GetOrder{
		Id: orderId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	req, err := stream.Recv()
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	meta := req.GetMetadata()
	fileData := bytes.Buffer{}

	for {
		logger.Debug("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			logger.Debug("no more data")
			break
		}

		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
			return
		}

		chunk := req.GetFile().Content

		_, err = fileData.Write(chunk)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
			return
		}
	}

	f, err := os.Create(meta.Name)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to create file")
		return
	}
	f.Write(fileData.Bytes())
	defer os.Remove(meta.Name)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprintf("%d", meta.Size))
	c.Header("Content-Disposition", "attachment; filename="+meta.GetName())
	c.File(meta.Name)
}

func (h *OrderHandler) create(c *gin.Context) {
	userId, exists := c.Get(middleware.UserIdCtx)
	if !exists {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty user id param")
		return
	}

	var dto *order_api.CreateOrder
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	dto.UserId = userId.(string)

	// _, err := h.orderApi.Create(c, &order_api.CreateOrder{
	// 	Id:     dto.Id,
	// 	Count:  dto.Count,
	// 	UserId: dto.UserId,
	// 	//TODO дописать преобразование позиций
	// 	// Positions: dto.Positions,
	// })
	_, err := h.orderApi.Create(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	// c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/orders/%s", order.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *OrderHandler) save(c *gin.Context) {
	userId, exists := c.Get(middleware.UserIdCtx)
	if !exists {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty user id param")
		return
	}

	var dto *order_api.CreateOrder
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	dto.UserId = userId.(string)

	_, err := h.orderApi.Save(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	// c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/orders/%s", order.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Saved"})
}
