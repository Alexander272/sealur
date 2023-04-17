package new_pro

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/order_model"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderApi order_api.OrderServiceClient
	userApi  user_api.UserServiceClient
	emailApi email_api.EmailServiceClient
	fileApi  file_api.FileServiceClient
}

func NewOrderHandler(orderApi order_api.OrderServiceClient, userApi user_api.UserServiceClient, emailApi email_api.EmailServiceClient,
	fileApi file_api.FileServiceClient,
) *OrderHandler {
	return &OrderHandler{
		orderApi: orderApi,
		userApi:  userApi,
		emailApi: emailApi,
		fileApi:  fileApi,
	}
}

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	handler := NewOrderHandler(h.orderApi, h.userApi, h.emailApi, h.fileApi)

	order := api.Group("/orders", h.middleware.UserIdentity)
	{
		// order.GET("/:id", handler.get)
		order.GET("/current", handler.getCurrent)
		order.GET("/all", handler.getAll)
		// order.GET("/open", handler.getOpen)
		order.GET("/:id/заявка.zip", handler.getFile)
		order.POST("/", handler.create)
		order.POST("/copy", handler.copy)
		order.POST("/save", handler.save)
		manager := order.Group("/", h.middleware.AccessForManager)
		{
			manager.GET("/:id", handler.get)
			manager.GET("/open", handler.getOpen)
			manager.POST("/finish", handler.finish)
			manager.POST("/manager", handler.setManager)
		}
		// order.POST("/finish", handler.finish)
		// order.POST("/manager", handler.setManager)
	}
}

func (h *OrderHandler) get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty id param")
		return
	}

	order, err := h.orderApi.Get(c, &order_api.GetOrder{Id: id})
	if err != nil {
		if strings.Contains(err.Error(), "order is not exist") {
			models.NewErrorResponse(c, http.StatusBadRequest, "order is not exist", "Ошибка: Данный заказ не найден")
			return
		}
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		return
	}

	user, err := h.userApi.Get(c, &user_api.GetUser{Id: order.Order.UserId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		return
	}
	order.User = user

	c.JSON(http.StatusOK, models.DataResponse{Data: order})
}

func (h *OrderHandler) getCurrent(c *gin.Context) {
	userId, exists := c.Get(middleware.UserIdCtx)
	if !exists {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty user id param")
		return
	}

	user, err := h.userApi.Get(c, &user_api.GetUser{Id: userId.(string)})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	order, err := h.orderApi.GetCurrent(c, &order_api.GetCurrentOrder{UserId: user.Id, ManagerId: user.ManagerId})
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

func (h *OrderHandler) getOpen(c *gin.Context) {
	userId, exists := c.Get(middleware.UserIdCtx)
	if !exists {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty user id param")
		return
	}

	orders, err := h.orderApi.GetOpen(c, &order_api.GetManagerOrders{ManagerId: userId.(string)})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
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
	// дописать преобразование позиций
	// Positions: dto.Positions,
	// })
	_, err := h.orderApi.Create(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	// c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/orders/%s", order.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Created"})
}

func (h *OrderHandler) copy(c *gin.Context) {
	var dto *order_api.CopyOrder
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "invalid data send")
		return
	}

	_, err := h.orderApi.Copy(c, dto)
	if err != nil {
		if strings.Contains(err.Error(), "position exists") {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Ошибка: Одна или несколько позиций дублируются")
			return
		}
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "произошла ошибка во время копирования")
		return
	}

	_, err = h.fileApi.CopyGroup(c, &file_api.CopyGroupRequest{
		Bucket:   "pro",
		Group:    dto.FromId,
		NewGroup: dto.TargetId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "произошла ошибка во время копирования чертежей")
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Copied successfully"})
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

	data, err := h.userApi.GetManagerEmail(c, &user_api.GetUser{Id: dto.UserId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}

	_, err = h.emailApi.SendNotification(c, &email_api.NotificationData{Email: data.Email, User: data.User, OrderId: dto.Id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to send email")
		return
	}

	// c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/orders/%s", order.Id))
	c.JSON(http.StatusOK, models.IdResponse{Message: "Saved"})
}

func (h *OrderHandler) finish(c *gin.Context) {
	var dto *order_api.Status
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "отправлены некорректные данные")
		return
	}
	dto.Status = order_model.OrderStatus_finish

	_, err := h.orderApi.SetStatus(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated status successfully"})
}

func (h *OrderHandler) setManager(c *gin.Context) {
	var dto *order_api.Manager
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "отправлены некорректные данные")
		return
	}

	//TODO отправлять email при изменении (? а надо ли это вообще)
	_, err := h.orderApi.SetManager(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated manager successfully"})
}
