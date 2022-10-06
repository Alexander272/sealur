package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/material_api"
)

func (h *MaterialsHandlers) CreateAlpha(ctx context.Context, alpha *material_api.CreateAlphaRequest) (*material_api.Response, error) {
	err := h.service.CreateAlpha(ctx, alpha)
	if err != nil {
		return nil, err
	}

	return &material_api.Response{}, nil
}

func (h *MaterialsHandlers) UpdateAlpha(ctx context.Context, alpha *material_api.UpdateAlphaRequest) (*material_api.Response, error) {
	if err := h.service.UpateAlpha(ctx, alpha); err != nil {
		return nil, err
	}
	return &material_api.Response{}, nil
}

func (h *MaterialsHandlers) DeleteAlpha(ctx context.Context, alpha *material_api.DeleteAlphaRequest) (*material_api.Response, error) {
	if err := h.service.DeleteAlpha(ctx, alpha); err != nil {
		return nil, err
	}
	return &material_api.Response{}, nil
}
