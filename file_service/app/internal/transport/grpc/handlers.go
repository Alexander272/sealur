package grpc

import (
	"context"

	"github.com/Alexander272/sealur/file_service/internal/config"
	"github.com/Alexander272/sealur/file_service/internal/service"
	proto_file "github.com/Alexander272/sealur/file_service/internal/transport/grpc/proto"
)

type Handler struct {
	service *service.Service
	conf    config.ApiConfig
}

func NewHandler(service *service.Service, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,
	}
}

func (h *Handler) Ping(ctx context.Context, req *proto_file.PingRequest) (*proto_file.PingResponse, error) {
	return &proto_file.PingResponse{Ping: "pong"}, nil
}
