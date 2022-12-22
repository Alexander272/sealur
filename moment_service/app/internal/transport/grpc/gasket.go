package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
)

type GasketHandlers struct {
	service service.Gasket
	gasket_api.UnimplementedGasketServiceServer
}

func NewGasketService(service service.Gasket) *GasketHandlers {
	return &GasketHandlers{service: service}
}

func (h *GasketHandlers) GetFullData(ctx context.Context, req *gasket_api.GetFullDataRequest) (*gasket_api.FullDataResponse, error) {
	data, err := h.service.GetData(ctx, req)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (h *GasketHandlers) GetGasket(ctx context.Context, req *gasket_api.GetGasketRequest) (*gasket_api.GasketResponse, error) {
	gasket, err := h.service.GetGasket(ctx, req)
	if err != nil {
		return nil, err
	}

	return &gasket_api.GasketResponse{Gasket: gasket}, nil
}

func (h *GasketHandlers) CreateGasket(ctx context.Context, gasket *gasket_api.CreateGasketRequest) (*gasket_api.IdResponse, error) {
	id, err := h.service.CreateGasket(ctx, gasket)
	if err != nil {
		return nil, err
	}

	return &gasket_api.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateGasket(ctx context.Context, gasket *gasket_api.UpdateGasketRequest) (*gasket_api.Response, error) {
	if err := h.service.UpdateGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) DeleteGasket(ctx context.Context, gasket *gasket_api.DeleteGasketRequest) (*gasket_api.Response, error) {
	if err := h.service.DeleteGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

//---
func (h *GasketHandlers) CreateManyGasketData(ctx context.Context, data *gasket_api.CreateManyGasketDataRequest) (*gasket_api.Response, error) {
	if err := h.service.CreateManyGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) CreateGasketData(ctx context.Context, data *gasket_api.CreateGasketDataRequest) (*gasket_api.Response, error) {
	if err := h.service.CreateGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) UpdateGasketData(ctx context.Context, data *gasket_api.UpdateGasketDataRequest) (*gasket_api.Response, error) {
	if err := h.service.UpdateGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) UpdateGasketTypeId(ctx context.Context, data *gasket_api.UpdateGasketTypeIdRequest) (*gasket_api.Response, error) {
	if err := h.service.UpdateGasketTypeId(ctx, data); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}

func (h *GasketHandlers) DeleteGasketData(ctx context.Context, data *gasket_api.DeleteGasketDataRequest) (*gasket_api.Response, error) {
	if err := h.service.DeleteGasketData(ctx, data); err != nil {
		return nil, err
	}
	return &gasket_api.Response{}, nil
}
