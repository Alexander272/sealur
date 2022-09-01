package grpc

import (
	"github.com/Alexander272/sealur/moment_service/internal/config"
	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type Ping interface {
	moment_api.PingServiceServer
}

type CalcFlange interface {
	moment_api.CalcFlangeServiceServer
}

type CalcCap interface {
	moment_api.CalcCapServiceServer
}

type Materials interface {
	moment_api.MaterialsServiceServer
}

type Gasket interface {
	moment_api.GasketServiceServer
}

type Flange interface {
	moment_api.FlangeServiceServer
}

type Read interface {
	moment_api.ReadServiceServer
}

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	Ping
	CalcFlange
	CalcCap
	Materials
	Gasket
	Flange
	Read
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service:    service,
		conf:       conf,
		Ping:       NewPingHandlers(),
		CalcFlange: NewCalcFlangeHandlers(service.CalcFlange),
		CalcCap:    NewCalcCapHandlers(service.CalcCap),
		Materials:  NewMaterialsHandlers(service.Materials),
		Gasket:     NewGasketService(service.Gasket),
		Flange:     NewFlangeHandlers(service.Flange),
		Read:       NewReadHandlers(service.Read),
	}
}
