package grpc

import (
	"context"

	"github.com/Alexander272/sealur/user_service/internal/config"
	"github.com/Alexander272/sealur/user_service/internal/service"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
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

func (h *Handler) Ping(ctx context.Context, req *proto_user.PingRequest) (*proto_user.PingResponse, error) {
	return &proto_user.PingResponse{Ping: "pong"}, nil
}
