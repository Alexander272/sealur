package file

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	fileClient file_api.FileServiceClient
	middleware *middleware.Middleware
}

func NewHandler(middleware *middleware.Middleware) *Handler {
	return &Handler{middleware: middleware}
}

func (h *Handler) InitRoutes(conf config.ServicesConfig, api *gin.RouterGroup) {
	//* file service connect
	//TODO стоит ли так оставлять сертификат?
	//* определение сертификата
	// creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	// if err != nil {
	// 	logger.Fatalf("failed to load certificate. error: %w", err)
	// }

	//* данные для аутентификации
	// auth := models.Authentication{
	// 	ServiceName: conf.FileService.AuthName,
	// 	Password:    conf.FileService.AuthPassword,
	// }

	//* опции grpc
	// opts := []grpc.DialOption{
	// 	grpc.WithTransportCredentials(creds),
	// 	grpc.WithPerRPCCredentials(&auth),
	// }

	//* подключение к сервису
	// connect, err := grpc.Dial(conf.FileService.Url, opts...)
	connect, err := grpc.Dial(conf.FileService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed connection to file service. error: %s", err.Error())
	}

	fileClient := file_api.NewFileServiceClient(connect)
	h.fileClient = fileClient

	files := api.Group("/files")
	{
		files.GET("/ping", h.pingPro)

		h.initFilesRoutes(files)
	}
}

func (h *Handler) pingPro(c *gin.Context) {
	res, err := h.fileClient.Ping(c, &file_api.PingRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: res.Ping})
}
