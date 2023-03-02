package new_pro

import (
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro/snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_standard_api"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Handler struct {
	middleware     *middleware.Middleware
	snpApi         snp_api.SnpDataServiceClient
	snpStandardApi snp_standard_api.SnpStandardServiceClient
	// pingClient      moment.PingServiceClient
	// gasketClient    gasket_api.GasketServiceClient
	// materialsClient material_api.MaterialsServiceClient
	// flangeClient    flange_api.FlangeServiceClient
	// deviceClient    device_api.DeviceServiceClient
	// readClient      read_api.ReadServiceClient
	// calcClient      calc_api.CalcServiceClient
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
		ServiceName: conf.ProService.AuthName,
		Password:    conf.ProService.AuthPassword,
	}

	//* опции grpc
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&auth),
	}

	//* подключение к сервису
	connect, err := grpc.Dial(conf.ProService.Url, opts...)
	if err != nil {
		logger.Fatalf("failed connection to pro service. error: %w", err)
	}

	h.snpApi = snp_api.NewSnpDataServiceClient(connect)
	h.snpStandardApi = snp_standard_api.NewSnpStandardServiceClient(connect)

	pro := api.Group("/sealur-pro")
	{
		// pro.GET("/ping", h.ping)
	}
	h.initSNPRoutes(pro)
	h.initSnpStandardRoutes(pro)
}

// func (h *Handler) ping(c *gin.Context) {
// 	res, err := h.pingClient.Ping(c, &moment.PingRequest{})
// 	if err != nil {
// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}
// 	c.JSON(http.StatusOK, models.DataResponse{Data: res.Ping})
// }
