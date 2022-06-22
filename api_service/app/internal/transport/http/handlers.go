package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Alexander272/sealur/api_service/docs"
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/service"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	httpV1 "github.com/Alexander272/sealur/api_service/internal/transport/http/v1"
	"github.com/Alexander272/sealur/api_service/pkg/limiter"
	"github.com/gin-contrib/cors"
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
			AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5000", "http://localhost:8080"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"access-control-allow-origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "authorization", "accept", "origin", "Origin", "Cache-Control", "X-Requested-With"},
			ExposeHeaders:    []string{"access-control-allow-origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "authorization", "accept", "origin", "Origin", "Cache-Control", "X-Requested-With"},
			AllowCredentials: true,
			// AllowOriginFunc: func(origin string) bool {
			// 	logger.Debug(origin)
			// 	return origin == "http://localhost:3000"
			// },
			MaxAge: 12 * time.Hour,
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

	h.initAPI(conf.Services, conf.Auth, router)

	return router
}

func (h *Handler) initAPI(conf config.ServicesConfig, auth config.AuthConfig, router *gin.Engine) {
	handlerV1 := httpV1.NewHandler(h.services, middleware.NewMiddleware(h.services, auth))
	api := router.Group("/api")
	{
		handlerV1.Init(conf, auth, api)
	}
}
