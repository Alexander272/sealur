package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/config"
	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	pro_api.UnimplementedProServiceServer
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,
	}
}

func (h *Handler) Ping(ctx context.Context, req *pro_api.PingRequest) (*pro_api.PingResponse, error) {
	return &pro_api.PingResponse{Ping: "pong"}, nil
}
