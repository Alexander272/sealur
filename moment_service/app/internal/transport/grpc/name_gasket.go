package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (h *DeviceHandlers) GetNameGasket(ctx context.Context, req *device_api.GetNameGasketRequest) (*device_api.NameGasketResponse, error) {
	gasket, err := h.service.GetNameGasket(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.NameGasketResponse{Gasket: gasket}, nil
}

func (h *DeviceHandlers) CreateNameGasket(ctx context.Context, gasket *device_api.CreateNameGasketRequest) (*device_api.IdResponse, error) {
	id, err := h.service.CreateNameGasket(ctx, gasket)
	if err != nil {
		return nil, err
	}
	return &device_api.IdResponse{Id: id}, nil
}

func (h *DeviceHandlers) CreateFewNameGasket(ctx context.Context, gasket *device_api.CreateFewNameGasketRequest) (*device_api.Response, error) {
	if err := h.service.CreateFewNameGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) UpdateNameGasket(ctx context.Context, gasket *device_api.UpdateNameGasketRequest) (*device_api.Response, error) {
	if err := h.service.UpdateNameGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) DeleteNameGsket(ctx context.Context, gasket *device_api.DeleteNameGasketRequest) (*device_api.Response, error) {
	if err := h.service.DeleteNameGsket(ctx, gasket); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
