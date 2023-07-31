package new_pro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
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
	botApi   api.MostBotApi
}

func NewOrderHandler(orderApi order_api.OrderServiceClient, userApi user_api.UserServiceClient, emailApi email_api.EmailServiceClient,
	fileApi file_api.FileServiceClient, botApi api.MostBotApi,
) *OrderHandler {
	return &OrderHandler{
		orderApi: orderApi,
		userApi:  userApi,
		emailApi: emailApi,
		fileApi:  fileApi,
		botApi:   botApi,
	}
}

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	handler := NewOrderHandler(h.orderApi, h.userApi, h.emailApi, h.fileApi, h.botApi)

	order := api.Group("/orders", h.middleware.UserIdentity)
	{
		// order.GET("/:id", handler.get)
		order.GET("/current", handler.getCurrent)
		order.GET("/all", handler.getAll)
		// order.GET("/open", handler.getOpen)
		order.GET("/:id/заявка.zip", handler.getFile)
		order.GET("/analytics", handler.getAnalytics)
		order.GET("/analytics/count", handler.getCountAnalytics)
		order.GET("/analytics/full", handler.getFullAnalytics)
		order.POST("/", handler.create)
		order.POST("/copy", handler.copy)
		order.POST("/save", handler.save)
		order.PUT("/info", handler.setInfo)
		manager := order.Group("/", h.middleware.AccessForManager)
		{
			manager.GET("/:id", handler.get)
			manager.GET("/open", handler.getOpen)
			manager.POST("/finish", handler.finish)
			manager.POST("/manager", handler.setManager)

			manager.GET("/last", handler.getLast)
			manager.GET("/number/:number", handler.getByNumber)
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
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "id": "%s" }`, id))
		return
	}

	user, err := h.userApi.Get(c, &user_api.GetUser{Id: order.Order.UserId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "userId": "%s" }`, order.Order.UserId))
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
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "userId": "%s" }`, userId))
		return
	}

	order, err := h.orderApi.GetCurrent(c, &order_api.GetCurrentOrder{UserId: user.Id, ManagerId: user.ManagerId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "userId": "%s", "managerId": "%s" }`, user.Id, user.ManagerId))
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
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "userId": "%s" }`, userId))
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
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "managerId": "%s" }`, userId))
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
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "orderId": "%s" }`, orderId))
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
			h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "orderId": "%s" }`, orderId))
			return
		}

		chunk := req.GetFile().Content

		_, err = fileData.Write(chunk)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
			h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "orderId": "%s" }`, orderId))
			return
		}
	}

	f, err := os.Create(meta.Name)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to create file")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "orderId": "%s" }`, orderId))
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

func (h *OrderHandler) getByNumber(c *gin.Context) {
	number := c.Param("number")

	order, err := h.orderApi.GetByNumber(c, &order_api.GetOrderByNumber{Number: number})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "number": "%s" }`, number))
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: order})
}

func (h *OrderHandler) getLast(c *gin.Context) {
	orders, err := h.orderApi.GetLast(c, &order_api.GetLastOrders{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		h.botApi.SendError(c, err.Error(), "")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: orders})
}

func (h *OrderHandler) getAnalytics(c *gin.Context) {
	periodAt := c.Query("periodAt")
	if periodAt == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty user id param")
		return
	}
	periodEnd := c.Query("periodEnd")
	if periodEnd == "" {
		models.NewErrorResponse(c, http.StatusBadRequest, "empty param", "empty user id param")
		return
	}

	analytics, err := h.orderApi.GetAnalytics(c, &order_api.GetOrderAnalytics{PeriodAt: periodAt, PeriodEnd: periodEnd})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "periodAt": "%s", "periodEnd": "%s" }`, periodAt, periodEnd))
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: analytics})
}

func (h *OrderHandler) getCountAnalytics(c *gin.Context) {
	orders, err := h.orderApi.GetOrderCount(c, &order_api.GetOrderCountAnalytics{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		h.botApi.SendError(c, err.Error(), "")
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: orders.OrderCount})
}

func (h *OrderHandler) getFullAnalytics(c *gin.Context) {
	periodAt := c.Query("periodAt")
	periodEnd := c.Query("periodEnd")
	userId := c.Query("userId")

	data, err := h.orderApi.GetBidAnalytics(c, &order_api.GetFullOrderAnalytics{PeriodAt: periodAt, PeriodEnd: periodEnd, UserId: userId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "periodAt": "%s", "periodEnd": "%s", "userId": "%s" }`, periodAt, periodEnd, userId))
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: data.Orders})
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

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

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

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}

	_, err = h.fileApi.CopyGroup(c, &file_api.CopyGroupRequest{
		Bucket:   "pro",
		Group:    dto.FromId,
		NewGroup: dto.TargetId,
	})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "произошла ошибка во время копирования чертежей")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

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

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}

	data, err := h.userApi.GetManagerEmail(c, &user_api.GetUser{Id: dto.UserId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "managerId": "%s" }`, userId))
		return
	}

	_, err = h.emailApi.SendNotification(c, &email_api.NotificationData{Email: data.Email, User: data.User, OrderId: dto.Id})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "failed to send email")
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "Email": "%s" }`, data.Email))
		return
	}

	// c.Header("Location", fmt.Sprintf("/api/v1/sealur-pro/orders/%s", order.Id))
	c.JSON(http.StatusCreated, models.IdResponse{Message: "Saved"})
}

func (h *OrderHandler) setInfo(c *gin.Context) {
	var dto *order_api.Info
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "отправлены некорректные данные")
		return
	}

	_, err := h.orderApi.SetInfo(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated"})
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

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

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

	user, err := h.userApi.Get(c, &user_api.GetUser{Id: dto.UserId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "userId": "%s" }`, dto.UserId))
		return
	}

	manager, err := h.userApi.Get(c, &user_api.GetUser{Id: dto.OldManagerId})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "oldManagerId": "%s" }`, dto.OldManagerId))
		return
	}

	// отправлять email при изменении (? а надо ли это вообще)
	//? если отдавать новому менеджеру фио того кто переслал и данные о заказчике, то мне надо сделать еще два запроса для получения этих данных
	_, err = h.orderApi.SetManager(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}

	_, err = h.emailApi.SendRedirect(c, &email_api.RedirectData{Email: dto.ManagerEmail, OrderId: dto.OrderId, User: user, Manager: manager})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось отправить email: "+err.Error())
		h.botApi.SendError(c, err.Error(), fmt.Sprintf(`{ "Email": "%s" }`, dto.ManagerEmail))
		return
	}

	c.JSON(http.StatusOK, models.IdResponse{Message: "Updated manager successfully"})
}
