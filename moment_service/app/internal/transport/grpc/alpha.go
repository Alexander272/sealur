package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (h *MaterialsHandlers) CreateAlpha(ctx context.Context, alpha *moment_api.CreateAlphaRequest) (*moment_api.Response, error) {
	err := h.service.CreateAlpha(ctx, alpha)
	if err != nil {
		return nil, err
	}

	return &moment_api.Response{}, nil
}

func (h *MaterialsHandlers) UpdateAlpha(ctx context.Context, alpha *moment_api.UpdateAlphaRequest) (*moment_api.Response, error) {
	if err := h.service.UpateAlpha(ctx, alpha); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *MaterialsHandlers) DeleteAlpha(ctx context.Context, alpha *moment_api.DeleteAlphaRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteAlpha(ctx, alpha); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}
