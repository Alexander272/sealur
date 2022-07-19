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

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	Ping
	Flange
	Materials
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service:   service,
		conf:      conf,
		Ping:      NewPingHandlers(),
		Flange:    NewFlangeHandlers(service.Flange),
		Materials: NewMaterialsHandlers(service.Materials),
	}
}
