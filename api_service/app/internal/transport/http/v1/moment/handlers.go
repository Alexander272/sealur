package moment

import (
	"net/http"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/moment"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/Alexander272/sealur_proto/api/moment/read_api"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Handler struct {
	middleware      *middleware.Middleware
	pingClient      moment.PingServiceClient
	gasketClient    gasket_api.GasketServiceClient
	materialsClient material_api.MaterialsServiceClient
	flangeClient    flange_api.FlangeServiceClient
	readClient      read_api.ReadServiceClient
	calcClient      calc_api.CalcServiceClient
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

	pingClient := moment.NewPingServiceClient(connect)
	gasketClient := gasket_api.NewGasketServiceClient(connect)
	materialsClient := material_api.NewMaterialsServiceClient(connect)
	flangeClient := flange_api.NewFlangeServiceClient(connect)
	readClient := read_api.NewReadServiceClient(connect)
	calcClient := calc_api.NewCalcServiceClient(connect)

	h.pingClient = pingClient
	h.gasketClient = gasketClient
	h.materialsClient = materialsClient
	h.flangeClient = flangeClient
	h.readClient = readClient
	h.calcClient = calcClient

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

		h.initReadRoutes(moment)

		h.initCalcRoutes(moment)
	}
}

func (h *Handler) pingUsers(c *gin.Context) {
	res, err := h.pingClient.Ping(c, &moment.PingRequest{})
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
		return
	}
	c.JSON(http.StatusOK, models.DataResponse{Data: res.Ping})
}
