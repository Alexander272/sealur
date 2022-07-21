package grpc

import "github.com/Alexander272/sealur/moment_service/internal/service"

type FlangeHandlers struct {
	service service.Flange
}

func NewFlangeHandlers(service service.Flange) *FlangeHandlers {
	return &FlangeHandlers{
		service: service,
	}
}
