package user

import (
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/service"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	userApi    user_api.UserServiceClient
	emailApi   email_api.EmailServiceClient
	conf       *config.Config
	services   *service.Services
	middleware *middleware.Middleware
	cookieName string
}

func NewHandler(conf *config.Config, services *service.Services, middleware *middleware.Middleware, cookieName string) *Handler {
	return &Handler{
		conf:       conf,
		services:   services,
		middleware: middleware,
		cookieName: cookieName,
	}
}

func (h *Handler) InitRoutes(conf config.ServicesConfig, api *gin.RouterGroup) {
	//* user service connect

	//* определение сертификата
	// creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	// if err != nil {
	// 	logger.Fatalf("failed to load certificate. error: %w", err)
	// }

	//* данные для аутентификации
	// auth := models.Authentication{
	// 	ServiceName: conf.UserService.AuthName,
	// 	Password:    conf.UserService.AuthPassword,
	// }

	//* опции grpc
	// opts := []grpc.DialOption{
	// 	grpc.WithTransportCredentials(creds),
	// 	grpc.WithPerRPCCredentials(&auth),
	// }

	//* подключение к сервису
	// connect, err := grpc.Dial(conf.UserService.Url, opts...)
	connect, err := grpc.Dial(conf.UserService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed connection to user service. error: %s", err.Error())
	}

	// userClient := user_api.NewUserServiceClient(connect)
	// h.userClient = userClient
	h.userApi = user_api.NewUserServiceClient(connect)

	//* данные для аутентификации
	// authEmail := models.Authentication{
	// 	ServiceName: conf.EmailService.AuthName,
	// 	Password:    conf.EmailService.AuthPassword,
	// }
	//* опции grpc
	// optsEmail := []grpc.DialOption{
	// 	grpc.WithTransportCredentials(creds),
	// 	grpc.WithPerRPCCredentials(&authEmail),
	// }
	//* подключение к сервису
	// connectEmail, err := grpc.Dial(conf.EmailService.Url, optsEmail...)
	connectEmail, err := grpc.Dial(conf.EmailService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed connection to email service. error: %s", err.Error())
	}

	// emailClient := email_api.NewEmailServiceClient(connectEmail)
	// h.emailClient = emailClient
	h.emailApi = email_api.NewEmailServiceClient(connectEmail)

	root := api.Group("/")
	{
		// users.GET("/users/ping", h.pingUsers)

		// h.initUserRoutes(users)
	}
	h.initAuthRoutes(root)
	h.initUserRoutes(root)
}

// func (h *Handler) pingUsers(c *gin.Context) {
// 	res, err := h.userClient.Ping(c, &user_api.PingRequest{})
// 	if err != nil {
// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}
// 	c.JSON(http.StatusOK, models.DataResponse{Data: res.Ping})
// }
