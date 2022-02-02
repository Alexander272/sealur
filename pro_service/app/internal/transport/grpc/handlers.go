package grpc

import "github.com/Alexander272/sealur/pro_service/internal/service"

type Handler struct {
	service *service.Services
}

func NewHandler(service *service.Services) *Handler {
	return &Handler{
		service: service,
	}
}
