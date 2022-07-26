package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (h *GasketHandlers) GetGasketType(ctx context.Context, req *moment_api.GetGasketTypeRequest) (*moment_api.GasketTypeResponse, error) {
	typeGasket, err := h.service.GetTypeGasket(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_api.GasketTypeResponse{GasketType: typeGasket}, nil
}

func (h *GasketHandlers) CreateGasketType(ctx context.Context, gasket *moment_api.CreateGasketTypeRequest) (*moment_api.IdResponse, error) {
	id, err := h.service.CreateTypeGasket(ctx, gasket)
	if err != nil {
		return nil, err
	}

	return &moment_api.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateGasketType(ctx context.Context, gasket *moment_api.UpdateGasketTypeRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateTypeGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *GasketHandlers) DeleteGasketType(ctx context.Context, gasket *moment_api.DeleteGasketTypeRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteTypeGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}
