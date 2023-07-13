package new_pro

import (
	"encoding/json"
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/api"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/gin-gonic/gin"
)

type ConnectHandler struct {
	emailApi email_api.EmailServiceClient
	botApi   api.MostBotApi
}

func NewConnectHandler(emailApi email_api.EmailServiceClient, botApi api.MostBotApi) *ConnectHandler {
	return &ConnectHandler{
		emailApi: emailApi,
		botApi:   botApi,
	}
}

func (h *Handler) initConnectRoutes(api *gin.RouterGroup) {
	handler := NewConnectHandler(h.emailApi, h.botApi)

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

		body, bodyErr := json.Marshal(dto)
		if bodyErr != nil {
			logger.Error("body error: ", bodyErr)
		}
		h.botApi.SendError(c, err.Error(), string(body))

		return
	}
	c.JSON(http.StatusOK, models.IdResponse{Message: "Sended feedback successfully"})
}
