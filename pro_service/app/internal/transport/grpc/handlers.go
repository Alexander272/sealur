package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/config"
	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,
	}
}

func (h *Handler) Ping(ctx context.Context, req *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{Ping: "pong"}, nil
}
