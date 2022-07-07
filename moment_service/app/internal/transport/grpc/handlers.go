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

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	Ping
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,
		Ping:    NewPingHandlers(),
	}
}
