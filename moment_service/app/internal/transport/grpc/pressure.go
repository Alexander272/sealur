package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (h *DeviceHandlers) GetPressure(ctx context.Context, req *device_api.GetPressureRequest) (*device_api.PressureResponse, error) {
	pressure, err := h.service.GetPressure(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.PressureResponse{Pressures: pressure}, nil
}

func (h *DeviceHandlers) CreatePressure(ctx context.Context, pres *device_api.CreatePressureRequest) (*device_api.IdResponse, error) {
	id, err := h.service.CreatePressure(ctx, pres)
	if err != nil {
		return nil, err
	}
	return &device_api.IdResponse{Id: id}, nil
}

func (h *DeviceHandlers) CreateFewPressure(ctx context.Context, pres *device_api.CreateFewPressureRequest) (*device_api.Response, error) {
	if err := h.service.CreateFewPressure(ctx, pres); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) UpdatePressure(ctx context.Context, pres *device_api.UpdatePressureRequest) (*device_api.Response, error) {
	if err := h.service.UpdatePressure(ctx, pres); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) DeletePressure(ctx context.Context, pres *device_api.DeletePressureRequest) (*device_api.Response, error) {
	if err := h.service.DeletePressure(ctx, pres); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
