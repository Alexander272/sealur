package models

import (
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type DataResponse struct {
	Data  interface{} `json:"data"`
	Count int         `json:"count,omitempty"`
}

type IdResponse struct {
	Id      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type FileResponse struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	OrigName string `json:"origName"`
	Link     string `json:"link"`
	Group    string `json:"group"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type ZipResponse struct {
	Name string `json:"name"`
	File []byte `json:"file"`
}

func NewErrorResponse(c *gin.Context, statusCode int, err, message string) {
	logger.Errorf("Url: %s | ClientIp: %s | ErrorResponse: %s", c.Request.URL, c.ClientIP(), err)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}
