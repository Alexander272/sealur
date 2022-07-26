package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type GasketHandlers struct {
	service service.Gasket
	moment_api.UnimplementedGasketServiceServer
}

func NewGasketService(service service.Gasket) *GasketHandlers {
	return &GasketHandlers{service: service}
}

func (h *GasketHandlers) GetGasket(ctx context.Context, req *moment_api.GetGasketRequest) (*moment_api.GasketResponse, error) {
	gasket, err := h.service.GetGasket(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_api.GasketResponse{Gasket: gasket}, nil
}

func (h *GasketHandlers) CreateGasket(ctx context.Context, gasket *moment_api.CreateGasketRequest) (*moment_api.IdResponse, error) {
	id, err := h.service.CreateGasket(ctx, gasket)
	if err != nil {
		return nil, err
	}

	return &moment_api.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateGasket(ctx context.Context, gasket *moment_api.UpdateGasketRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *GasketHandlers) DeleteGasket(ctx context.Context, gasket *moment_api.DeleteGasketRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

//---

func (h *GasketHandlers) CreateGasketData(ctx context.Context, data *moment_api.CreateGasketDataRequest) (*moment_api.Response, error) {
	if err := h.service.CreateGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *GasketHandlers) UpdateGasketData(ctx context.Context, data *moment_api.UpdateGasketDataRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *GasketHandlers) DeleteGasketData(ctx context.Context, data *moment_api.DeleteGasketDataRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}
