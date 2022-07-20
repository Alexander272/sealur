package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/config"
	"github.com/Alexander272/sealur/moment_service/internal/service"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type Ping interface {
	Ping(context.Context, *moment_proto.PingRequest) (*moment_proto.PingResponse, error)
}

type Flange interface {
	CalculateFlange(context.Context, *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error)
}

type Materials interface {
	GetMaterials(context.Context, *moment_proto.GetMaterialsRequest) (*moment_proto.MaterialsResponse, error)
	GetMaterialsWithIsEmpty(context.Context, *moment_proto.GetMaterialsRequest) (*moment_proto.MaterialsWithIsEmptyResponse, error)
	GetMaterialsData(context.Context, *moment_proto.GetMaterialsDataRequest) (*moment_proto.MaterialsDataResponse, error)
	CreateMaterial(context.Context, *moment_proto.CreateMaterialRequest) (*moment_proto.IdResponse, error)
	UpdateMaterial(context.Context, *moment_proto.UpdateMaterialRequest) (*moment_proto.Response, error)
	DeleteMaterial(context.Context, *moment_proto.DeleteMaterialRequest) (*moment_proto.Response, error)

	CreateVoltage(context.Context, *moment_proto.CreateVoltageRequest) (*moment_proto.Response, error)
	UpdateVoltage(context.Context, *moment_proto.UpdateVoltageRequest) (*moment_proto.Response, error)
	DeleteVoltage(context.Context, *moment_proto.DeleteVoltageRequest) (*moment_proto.Response, error)

	CreateElasticity(context.Context, *moment_proto.CreateElasticityRequest) (*moment_proto.Response, error)
	UpdateElasticity(context.Context, *moment_proto.UpdateElasticityRequest) (*moment_proto.Response, error)
	DeleteElasticity(context.Context, *moment_proto.DeleteElasticityRequest) (*moment_proto.Response, error)

	CreateAlpha(context.Context, *moment_proto.CreateAlphaRequest) (*moment_proto.Response, error)
	UpdateAlpha(context.Context, *moment_proto.UpdateAlphaRequest) (*moment_proto.Response, error)
	DeleteAlpha(context.Context, *moment_proto.DeleteAlphaRequest) (*moment_proto.Response, error)
}

type Gasket interface {
	GetGasket(context.Context, *moment_proto.GetGasketRequest) (*moment_proto.GasketResponse, error)
	CreateGasket(context.Context, *moment_proto.CreateGasketRequest) (*moment_proto.IdResponse, error)
	UpdateGasket(context.Context, *moment_proto.UpdateGasketRequest) (*moment_proto.Response, error)
	DeleteGasket(context.Context, *moment_proto.DeleteGasketRequest) (*moment_proto.Response, error)

	GetGasketType(context.Context, *moment_proto.GetGasketTypeRequest) (*moment_proto.GasketTypeResponse, error)
	CreateGasketType(context.Context, *moment_proto.CreateGasketTypeRequest) (*moment_proto.IdResponse, error)
	UpdateGasketType(context.Context, *moment_proto.UpdateGasketTypeRequest) (*moment_proto.Response, error)
	DeleteGasketType(context.Context, *moment_proto.DeleteGasketTypeRequest) (*moment_proto.Response, error)

	GetEnv(context.Context, *moment_proto.GetEnvRequest) (*moment_proto.EnvResponse, error)
	CreateEnv(context.Context, *moment_proto.CreateEnvRequest) (*moment_proto.IdResponse, error)
	UpdateEnv(context.Context, *moment_proto.UpdateEnvRequest) (*moment_proto.Response, error)
	DeleteEnv(context.Context, *moment_proto.DeleteEnvRequest) (*moment_proto.Response, error)

	CreateEnvData(context.Context, *moment_proto.CreateEnvDataRequest) (*moment_proto.Response, error)
	UpdateEnvData(context.Context, *moment_proto.UpdateEnvDataRequest) (*moment_proto.Response, error)
	DeleteEnvData(context.Context, *moment_proto.DeleteEnvDataRequest) (*moment_proto.Response, error)

	CreateGasketData(context.Context, *moment_proto.CreateGasketDataRequest) (*moment_proto.Response, error)
	UpdateGasketData(context.Context, *moment_proto.UpdateGasketDataRequest) (*moment_proto.Response, error)
	DeleteGasketData(context.Context, *moment_proto.DeleteGasketDataRequest) (*moment_proto.Response, error)
}

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	Ping
	Flange
	Materials
	Gasket
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service:   service,
		conf:      conf,
		Ping:      NewPingHandlers(),
		Flange:    NewFlangeHandlers(service.Flange),
		Materials: NewMaterialsHandlers(service.Materials),
		Gasket:    NewGasketService(service.Gasket),
	}
}
