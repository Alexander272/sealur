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

func (h *DeviceHandlers) GetFullNameGasket(ctx context.Context, req *device_api.GetFullNameGasketRequest) (*device_api.FullNameGasketResponse, error) {
	gasket, err := h.service.GetFullNameGasket(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.FullNameGasketResponse{Gasket: gasket}, nil
}

//TODO возвращать нужно объект, а не массив
// func (h *DeviceHandlers) GetNameGasketSize(ctx context.Context, req *device_api.GetNameGasketSizeRequest) (*device_api.NameGasketSizeResponse, error) {
// 	gasket, err := h.service.GetNameGasketSize(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &device_api.NameGasketSizeResponse{Gasket: gasket}, nil
// }

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

func (h *DeviceHandlers) DeleteNameGasket(ctx context.Context, gasket *device_api.DeleteNameGasketRequest) (*device_api.Response, error) {
	if err := h.service.DeleteNameGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
