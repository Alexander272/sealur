package http

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/sealur/api_service/docs"
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/service"
	httpV1 "github.com/Alexander272/sealur/api_service/internal/transport/http/v1"
	"github.com/Alexander272/sealur/api_service/pkg/limiter"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init(conf *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		limiter.Limit(conf.Limiter.RPS, conf.Limiter.Burst, conf.Limiter.TTL),
		cors.New(cors.Config{
			AllowedOrigins:   []string{conf.Http.Host},
			AllowedMethods:   []string{"GET"},
			AllowedHeaders:   []string{"Origin"},
			ExposedHeaders:   []string{"Content-Length"},
			AllowCredentials: true,
		}),
	)

	docs.SwaggerInfo_swagger.Host = fmt.Sprintf("%s:%s", conf.Http.Host, conf.Http.Port)
	if conf.Environment != "dev" {
		docs.SwaggerInfo_swagger.Host = conf.Http.Host
	}

	if conf.Environment != "prod" {
		router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Init router
	router.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(conf.Services, router)

	return router
}

func (h *Handler) initAPI(conf config.ServicesConfig, router *gin.Engine) {
	handlerV1 := httpV1.NewHandler(h.services)
	api := router.Group("/api")
	{
		handlerV1.Init(conf, api)
	}
}
