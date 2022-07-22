package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/moment_proto"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Handler struct {
	middleware       *middleware.Middleware
	pingClient       moment_proto.PingServiceClient
	gasketClient     moment_proto.GasketServiceClient
	materialsClient  moment_proto.MaterialsServiceClient
	flangeClient     moment_proto.FlangeServiceClient
	calcFlangeClient moment_proto.CalcFlangeServiceClient
}

func NewHandler(middleware *middleware.Middleware) *Handler {
	return &Handler{
		middleware: middleware,
	}
}

func (h *Handler) InitRoutes(conf config.ServicesConfig, api *gin.RouterGroup) {
	//* moment service connect

	//* определение сертификата
	creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	if err != nil {
		logger.Fatalf("failed to load certificate. error: %w", err)
	}

	//* данные для аутентификации
	auth := models.Authentication{
		ServiceName: conf.MomentService.AuthName,
		Password:    conf.MomentService.AuthPassword,
	}

	//* опции grpc
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&auth),
	}

	//* подключение к сервису
	connect, err := grpc.Dial(conf.MomentService.Url, opts...)
	if err != nil {
		logger.Fatalf("failed connection to pro service. error: %w", err)
	}

	pingClient := moment_proto.NewPingServiceClient(connect)
	gasketClient := moment_proto.NewGasketServiceClient(connect)
	materialsClient := moment_proto.NewMaterialsServiceClient(connect)
	flangeClient := moment_proto.NewFlangeServiceClient(connect)
	calcFlangeClient := moment_proto.NewCalcFlangeServiceClient(connect)

	h.pingClient = pingClient
	h.gasketClient = gasketClient
	h.materialsClient = materialsClient
	h.flangeClient = flangeClient
	h.calcFlangeClient = calcFlangeClient

	moment := api.Group("/sealur-moment")
	{
		moment.GET("/ping", h.pingUsers)

		h.initGasketRoutes(moment)
		h.initTypeGasketRoutes(moment)
		h.initEnvRoutes(moment)
		h.initEnvDataRoutes(moment)
		h.initGasketDataRoutes(moment)

		h.initMaterialsRoutes(moment)
		h.initVoltageRoutes(moment)
		h.initElasticityRoutes(moment)
		h.initAlphaRoutes(moment)

		h.initBoltsRoutes(moment)
		h.initTypeFlangeRoutes(moment)
		h.initStandartsRoutes(moment)
		h.initFlangeRoutes(moment)
	}
}

func (h *Handler) pingUsers(c *gin.Context) {
	res, err := h.pingClient.Ping(c, &moment_proto.PingRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: res.Ping})
}
