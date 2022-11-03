package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (h *DeviceHandlers) GetFinningFactor(ctx context.Context, req *device_api.GetFinningFactorRequest) (*device_api.FinningFactorResponse, error) {
	factor, err := h.service.GetFinningFactor(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.FinningFactorResponse{Finning: factor}, nil
}

func (h *DeviceHandlers) CreateFinningFactor(ctx context.Context, factor *device_api.CreateFinningFactorRequest) (*device_api.IdResponse, error) {
	id, err := h.service.CreateFinningFactor(ctx, factor)
	if err != nil {
		return nil, err
	}
	return &device_api.IdResponse{Id: id}, nil
}

func (h *DeviceHandlers) CreateFewFinningFactor(ctx context.Context, factor *device_api.CreateFewFinningFactorRequest) (*device_api.Response, error) {
	if err := h.service.CreateFewFinningFactor(ctx, factor); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) UpdateFinningFactor(ctx context.Context, factor *device_api.UpdateFinningFactorRequest) (*device_api.Response, error) {
	if err := h.service.UpdateFinningFactor(ctx, factor); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) DeleteFinningFactor(ctx context.Context, factor *device_api.DeleteFinningFactorRequest) (*device_api.Response, error) {
	if err := h.service.DeleteFinningFactor(ctx, factor); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
