package grpc

import (
	"context"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/service"
	proto_email "github.com/Alexander272/sealur/email_service/internal/transport/grpc/proto"
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

func (h *Handler) Ping(ctx context.Context, req *proto_email.PingRequest) (*proto_email.PingResponse, error) {
	return &proto_email.PingResponse{Ping: "pong"}, nil
}
