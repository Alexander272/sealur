package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
)

func (h *GasketHandlers) GetGasketType(ctx context.Context, req *gasket_api.GetGasketTypeRequest) (*gasket_api.GasketTypeResponse, error) {
	typeGasket, err := h.service.GetTypeGasket(ctx, req)
	if err != nil {
		return nil, err
	}

	return &gasket_api.GasketTypeResponse{GasketType: typeGasket}, nil
}

func (h *GasketHandlers) CreateGasketType(ctx context.Context, gasket *gasket_api.CreateGasketTypeRequest) (*gasket_api.IdResponse, error) {
	id, err := h.service.CreateTypeGasket(ctx, gasket)
	if err != nil {
		return nil, err
	}

	return &gasket_api.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateGasketType(ctx context.Context, gasket *gasket_api.UpdateGasketTypeRequest) (*gasket_api.Response, error) {
	if err := h.service.UpdateTypeGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) DeleteGasketType(ctx context.Context, gasket *gasket_api.DeleteGasketTypeRequest) (*gasket_api.Response, error) {
	if err := h.service.DeleteTypeGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}
