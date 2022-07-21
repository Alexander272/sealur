package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type GasketHandlers struct {
	service service.Gasket
}

func NewGasketService(service service.Gasket) *GasketHandlers {
	return &GasketHandlers{service: service}
}

func (h *GasketHandlers) GetGasket(ctx context.Context, req *moment_proto.GetGasketRequest) (*moment_proto.GasketResponse, error) {
	gasket, err := h.service.GetGasket(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.GasketResponse{Gasket: gasket}, nil
}

func (h *GasketHandlers) CreateGasket(ctx context.Context, gasket *moment_proto.CreateGasketRequest) (*moment_proto.IdResponse, error) {
	id, err := h.service.CreateGasket(ctx, gasket)
	if err != nil {
		return nil, err
	}

	return &moment_proto.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateGasket(ctx context.Context, gasket *moment_proto.UpdateGasketRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) DeleteGasket(ctx context.Context, gasket *moment_proto.DeleteGasketRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

//---

func (h *GasketHandlers) CreateGasketData(ctx context.Context, data *moment_proto.CreateGasketDataRequest) (*moment_proto.Response, error) {
	if err := h.service.CreateGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) UpdateGasketData(ctx context.Context, data *moment_proto.UpdateGasketDataRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) DeleteGasketData(ctx context.Context, data *moment_proto.DeleteGasketDataRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
