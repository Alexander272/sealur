package new_pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/gin-gonic/gin"
)

type ConnectHandler struct {
	emailApi email_api.EmailServiceClient
}

func NewConnectHandler(emailApi email_api.EmailServiceClient) *ConnectHandler {
	return &ConnectHandler{
		emailApi: emailApi,
	}
}

func (h *Handler) initConnectRoutes(api *gin.RouterGroup) {
	handler := NewConnectHandler(h.emailApi)

	connect := api.Group("/connect")
	{
		connect.POST("/feedback", handler.feedback)
	}
}

func (h *ConnectHandler) feedback(c *gin.Context) {
	var dto *email_api.Feedback
	if err := c.BindJSON(&dto); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "отправлены некорректные данные")
		return
	}

	_, err := h.emailApi.SendFeedback(c, dto)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Не удалось отправить сообщение")
		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Sended feedback successfully"})
}
