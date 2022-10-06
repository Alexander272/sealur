package grpc

import (
	"github.com/Alexander272/sealur/moment_service/internal/config"
	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
	"github.com/Alexander272/sealur_proto/api/moment/read_api"
)

type Ping interface {
	moment.PingServiceServer
}

type Calc interface {
	calc_api.CalcServiceServer
}

type Materials interface {
	material_api.MaterialsServiceServer
}

type Gasket interface {
	gasket_api.GasketServiceServer
}

type Flange interface {
	flange_api.FlangeServiceServer
}

type Read interface {
	read_api.ReadServiceServer
}

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	Ping
	Materials
	Gasket
	Flange
	Read
	Calc
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service:   service,
		conf:      conf,
		Ping:      NewPingHandlers(),
		Materials: NewMaterialsHandlers(service.Materials),
		Gasket:    NewGasketService(service.Gasket),
		Flange:    NewFlangeHandlers(service.Flange),
		Read:      NewReadHandlers(service.Read),
		Calc:      NewCalcHandlers(service.Calc),
	}
}
