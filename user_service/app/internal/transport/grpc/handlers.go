package grpc

import (
	"context"

	"github.com/Alexander272/sealur/user_service/internal/config"
	"github.com/Alexander272/sealur/user_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/user_api"
)

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	user_api.UnimplementedUserServiceServer
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,
	}
}

func (h *Handler) Ping(ctx context.Context, req *user_api.PingRequest) (*user_api.PingResponse, error) {
	return &user_api.PingResponse{Ping: "pong"}, nil
}
