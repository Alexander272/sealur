package new_pro

import (
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/middleware"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/Alexander272/sealur_proto/api/pro/putg_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_data_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_size_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/snp_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/temperature_api"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	middleware *middleware.Middleware

	snpApi         snp_api.SnpDataServiceClient
	snpStandardApi snp_standard_api.SnpStandardServiceClient
	orderApi       order_api.OrderServiceClient
	positionApi    position_api.PositionServiceClient
	materialApi    material_api.MaterialServiceClient
	snpMaterialApi snp_material_api.SnpMaterialServiceClient
	snpFillerApi   snp_filler_api.SnpFillerServiceClient
	temperatureApi temperature_api.TemperatureServiceClient
	flangeTypeApi  flange_type_api.FlangeTypeServiceClient
	snpTypeApi     snp_type_api.SnpTypeServiceClient
	sizeApi        snp_size_api.SnpSizeServiceClient
	snpDataApi     snp_data_api.SnpDataServiceClient

	putgApi putg_api.PutgDataServiceClient

	userApi  user_api.UserServiceClient
	emailApi email_api.EmailServiceClient
	fileApi  file_api.FileServiceClient
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
	// creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "Example-Root-CA")
	// if err != nil {
	// 	logger.Fatalf("failed to load certificate. error: %w", err)
	// }

	//* данные для аутентификации
	// auth := models.Authentication{
	// 	ServiceName: conf.ProService.AuthName,
	// 	Password:    conf.ProService.AuthPassword,
	// }

	//* опции grpc
	// opts := []grpc.DialOption{
	// 	// grpc.WithTransportCredentials(creds),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	grpc.WithPerRPCCredentials(&auth),
	// }

	userConnect, err := grpc.Dial(conf.UserService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed connection to user service. error: %w", err)
	}
	h.userApi = user_api.NewUserServiceClient(userConnect)

	emailConnect, err := grpc.Dial(conf.EmailService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed connection to email service. error: %w", err)
	}
	h.emailApi = email_api.NewEmailServiceClient(emailConnect)

	fileConnect, err := grpc.Dial(conf.FileService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed connection to file service. error: %w", err)
	}
	h.fileApi = file_api.NewFileServiceClient(fileConnect)

	//* подключение к сервису
	connect, err := grpc.Dial(conf.ProService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed connection to new pro service. error: %w", err)
	}

	h.snpApi = snp_api.NewSnpDataServiceClient(connect)
	h.snpStandardApi = snp_standard_api.NewSnpStandardServiceClient(connect)
	h.orderApi = order_api.NewOrderServiceClient(connect)
	h.positionApi = position_api.NewPositionServiceClient(connect)
	h.materialApi = material_api.NewMaterialServiceClient(connect)
	h.snpMaterialApi = snp_material_api.NewSnpMaterialServiceClient(connect)
	h.snpFillerApi = snp_filler_api.NewSnpFillerServiceClient(connect)
	h.temperatureApi = temperature_api.NewTemperatureServiceClient(connect)
	h.flangeTypeApi = flange_type_api.NewFlangeTypeServiceClient(connect)
	h.snpTypeApi = snp_type_api.NewSnpTypeServiceClient(connect)
	h.sizeApi = snp_size_api.NewSnpSizeServiceClient(connect)
	h.snpDataApi = snp_data_api.NewSnpDataServiceClient(connect)

	h.putgApi = putg_api.NewPutgDataServiceClient(connect)

	pro := api.Group("/sealur-pro")
	{
		// pro.GET("/ping", h.ping)
	}
	h.initSNPRoutes(pro)
	h.initSnpStandardRoutes(pro)
	h.initOrderRoutes(pro)
	h.initPositionRoutes(pro)
	h.initConnectRoutes(pro)
	h.initMaterialRoutes(pro)
	h.initSnpMaterialRoutes(pro)
	h.initSnpFillerRoutes(pro)
	h.initTemperatureRoutes(pro)
	h.initFlangeTypeRoutes(pro)
	h.initSnpTypeRoutes(pro)
	h.initSizeRoutes(pro)
	h.initSnpDataRoutes(pro)

	h.initPutgRoutes(pro)
}

// func (h *Handler) ping(c *gin.Context) {
// 	res, err := h.pingClient.Ping(c, &moment.PingRequest{})
// 	if err != nil {
// 		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "something went wrong")
// 		return
// 	}
// 	c.JSON(http.StatusOK, models.DataResponse{Data: res.Ping})
// }
