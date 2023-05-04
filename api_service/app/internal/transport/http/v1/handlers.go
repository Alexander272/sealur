package v1

import (
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/service"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/file"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/moment"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/new_pro"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services   *service.Services
	middleware *middleware.Middleware
}

const CookieName = "sealur_pro_session"

func NewHandler(services *service.Services, middleware *middleware.Middleware) *Handler {
	middleware.CookieName = CookieName

	return &Handler{
		services:   services,
		middleware: middleware,
	}
}

func (h *Handler) Init(conf *config.Config, api *gin.RouterGroup) {
	// proHandler := pro.NewHandler(h.middleware)
	momentHandler := moment.NewHandler(h.middleware)
	fileHandler := file.NewHandler(h.middleware)
	userHandler := user.NewHandler(conf, h.services, h.middleware, CookieName)

	v1 := api.Group("/v1")
	{
		// proHandler.InitRoutes(conf, v1)
		momentHandler.InitRoutes(conf.Services, v1)
		fileHandler.InitRoutes(conf.Services, v1)
		userHandler.InitRoutes(conf.Services, v1)

		v1.GET("/", h.notImplemented)
	}
	new_pro.NewHandler(h.middleware).InitRoutes(conf.Services, v1)
}

func (h *Handler) notImplemented(c *gin.Context) {}
