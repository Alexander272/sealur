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
	CalculateFlange(ctx context.Context, req *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error)
}

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	Ping
	Flange
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,
		Ping:    NewPingHandlers(),
		Flange:  NewFlangeHandlers(service.Flange),
	}
}
