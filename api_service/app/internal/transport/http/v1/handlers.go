package v1

import (
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/service"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/file"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/pro"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services   *service.Services
	middleware *middleware.Middleware
}

const CookieName = "session"

func NewHandler(services *service.Services, middleware *middleware.Middleware) *Handler {
	middleware.CookieName = CookieName

	return &Handler{
		services:   services,
		middleware: middleware,
	}
}

func (h *Handler) Init(conf config.ServicesConfig, auth config.AuthConfig, api *gin.RouterGroup) {
	proHandler := pro.NewHandler(h.middleware)
	fileHandler := file.NewHandler(h.middleware)
	userHandler := user.NewHandler(auth, h.services, h.middleware, CookieName)

	v1 := api.Group("/v1")
	{
		proHandler.InitRoutes(conf, v1)
		fileHandler.InitRoutes(conf, v1)
		userHandler.InitRoutes(conf, v1)

		v1.GET("/", h.notImplemented)
	}
}

func (h *Handler) notImplemented(c *gin.Context) {}
