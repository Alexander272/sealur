package grpc

import (
	"context"

	"github.com/Alexander272/sealur/file_service/internal/config"
	"github.com/Alexander272/sealur/file_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/file_api"
)

type Handler struct {
	service *service.Service
	conf    config.ApiConfig
	file_api.UnimplementedFileServiceServer
}

func NewHandler(service *service.Service, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,
	}
}

func (h *Handler) Ping(ctx context.Context, req *file_api.PingRequest) (*file_api.PingResponse, error) {
	return &file_api.PingResponse{Ping: "pong"}, nil
}
