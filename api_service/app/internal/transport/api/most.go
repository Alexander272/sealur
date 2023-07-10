package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Alexander272/sealur/api_service/internal/models/bot"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type MostApi struct {
	URL string
}

type MostBotApi interface {
	Ping() string
	SendError(c *gin.Context, err string, request string)
}

func NewMostApi(url string) *MostApi {
	return &MostApi{
		URL: "http://" + url,
	}
}

func (h *MostApi) Ping() string {
	return "Pong"
}

func (h *MostApi) SendError(c *gin.Context, e string, request string) {
	var user, company string

	name, exist := c.Get("userName")
	if exist {
		user = name.(string)
	}
	cm, exist := c.Get("company")
	if exist {
		company = cm.(string)
	}

	message := bot.Message{
		Service: bot.Service{
			Name: "SealurPro",
			Id:   "pro",
		},
		Data: bot.MessageData{
			Date:    time.Now().Format("02/01/2006 - 15:04:05"),
			IP:      c.ClientIP(),
			URL:     fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String()),
			Error:   e,
			User:    user,
			Company: company,
			Request: request,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(message)
	if err != nil {
		logger.Error("failed to read struct. error: %w", err)
	}

	_, err = http.Post(fmt.Sprintf("%s/api/v1/mattermost/send", h.URL), "application/json", &buf)
	if err != nil {
		logger.Error("failed to send error to bot. error: %w", err)
	}
}
