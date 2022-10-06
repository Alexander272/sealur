package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/material_api"
)

func (h *MaterialsHandlers) CreateVoltage(ctx context.Context, voltage *material_api.CreateVoltageRequest) (*material_api.Response, error) {
	err := h.service.CreateVoltage(ctx, voltage)
	if err != nil {
		return nil, err
	}

	return &material_api.Response{}, nil
}

func (h *MaterialsHandlers) UpdateVoltage(ctx context.Context, voltage *material_api.UpdateVoltageRequest) (*material_api.Response, error) {
	if err := h.service.UpdateVoltage(ctx, voltage); err != nil {
		return nil, err
	}
	return &material_api.Response{}, nil
}

func (h *MaterialsHandlers) DeleteVoltage(ctx context.Context, voltage *material_api.DeleteVoltageRequest) (*material_api.Response, error) {
	if err := h.service.DeleteVoltage(ctx, voltage); err != nil {
		return nil, err
	}
	return &material_api.Response{}, nil
}
