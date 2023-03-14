package pro

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	proClient  pro_api.ProServiceClient
	middleware *middleware.Middleware
}

func NewHandler(middleware *middleware.Middleware) *Handler {
	return &Handler{middleware: middleware}
}

func (h *Handler) InitRoutes(conf config.ServicesConfig, api *gin.RouterGroup) {
	//* pro service connect
	//TODO стоит ли так оставлять сертификат?
	//* определение сертификата
	// creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	// if err != nil {
	// 	logger.Fatalf("failed to load certificate. error: %w", err)
	// }

	//* данные для аутентификации
	auth := models.Authentication{
		ServiceName: conf.ProService.AuthName,
		Password:    conf.ProService.AuthPassword,
	}

	//* опции grpc
	opts := []grpc.DialOption{
		// grpc.WithTransportCredentials(creds),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&auth),
	}

	//* подключение к сервису
	connect, err := grpc.Dial(conf.ProService.Url, opts...)
	if err != nil {
		logger.Fatalf("failed connection to pro service. error: %w", err)
	}
	proClient := pro_api.NewProServiceClient(connect)
	h.proClient = proClient

	pro := api.Group("/sealur-pro")
	{
		pro.GET("/ping", h.pingPro)

		h.initStandRoutes(pro)
		h.initFlangeRoutes(pro)
		h.initStFlRoutes(pro)
		h.initTypeFlRoutes(pro)
		h.initAdditRoutes(pro)
		h.initSizeRoutes(pro)
		h.initSNPRoutes(pro)
		h.initPutgImageRoutes(pro)
		h.initPutgRoutes(pro)
		h.initPutgmImageRoutes(pro)
		h.initPutgmRoutes(pro)
		h.initMaterialsRoutes(pro)
		h.initBoltMaterialsRoutes(pro)
		h.initSizeIntRoutes(pro)

		h.initOrderRoutes(pro)
		h.initInterviewRoutes(pro)
	}
}

func (h *Handler) pingPro(c *gin.Context) {
	res, err := h.proClient.Ping(c, &pro_api.PingRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: res.Ping})
}
