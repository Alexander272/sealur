package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (h *MaterialsHandlers) CreateVoltage(ctx context.Context, voltage *moment_api.CreateVoltageRequest) (*moment_api.Response, error) {
	err := h.service.CreateVoltage(ctx, voltage)
	if err != nil {
		return nil, err
	}

	return &moment_api.Response{}, nil
}

func (h *MaterialsHandlers) UpdateVoltage(ctx context.Context, voltage *moment_api.UpdateVoltageRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateVoltage(ctx, voltage); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *MaterialsHandlers) DeleteVoltage(ctx context.Context, voltage *moment_api.DeleteVoltageRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteVoltage(ctx, voltage); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}
