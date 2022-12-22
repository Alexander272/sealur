package grpc

import (
	"context"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/email_api"
)

type Handler struct {
	service *service.Services
	conf    config.ApiConfig
	email_api.UnimplementedEmailServiceServer
}

func NewHandler(service *service.Services, conf config.ApiConfig) *Handler {
	return &Handler{
		service: service,
		conf:    conf,
	}
}

func (h *Handler) Ping(ctx context.Context, req *email_api.PingRequest) (*email_api.PingResponse, error) {
	return &email_api.PingResponse{Ping: "pong"}, nil
}

func (h *Handler) SendTest(ctx context.Context, req *email_api.SendTestRequest) (*email_api.SuccessResponse, error) {
	if err := h.service.Test.SendEmail(req); err != nil {
		return nil, err
	}
	return &email_api.SuccessResponse{Success: true}, nil
}
