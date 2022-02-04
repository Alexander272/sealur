package v1

import (
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/service"
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
	// pro service connect
	// _, err := grpc.Dial(conf.ProService.Url)
	// connect, err := grpc.Dial(conf.ProService.Url)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// proService = proto.NewProServiceClient(connect)

	v1 := api.Group("/v1")
	{
		h.initProRoutes(v1)
		v1.GET("/", h.notImplemented)
	}
}

func (h *Handler) notImplemented(c *gin.Context) {}
