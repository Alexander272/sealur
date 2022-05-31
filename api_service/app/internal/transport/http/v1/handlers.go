package v1

import (
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/service"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/file"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/pro"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(conf config.ServicesConfig, api *gin.RouterGroup) {
	proHandler := pro.NewHandler()
	fileHandler := file.NewHandler()

	v1 := api.Group("/v1")
	{
		proHandler.InitRoutes(conf, v1)
		fileHandler.InitRoutes(conf, v1)

		v1.GET("/", h.notImplemented)
	}
}

func (h *Handler) notImplemented(c *gin.Context) {}
