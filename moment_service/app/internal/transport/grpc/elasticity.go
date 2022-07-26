package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (h *MaterialsHandlers) CreateElasticity(ctx context.Context, elasticity *moment_api.CreateElasticityRequest) (*moment_api.Response, error) {
	err := h.service.CreateElasticity(ctx, elasticity)
	if err != nil {
		return nil, err
	}

	return &moment_api.Response{}, nil
}

func (h *MaterialsHandlers) UpdateElasticity(ctx context.Context, elasticity *moment_api.UpdateElasticityRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateElasticity(ctx, elasticity); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *MaterialsHandlers) DeleteElasticity(ctx context.Context, elasticity *moment_api.DeleteElasticityRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteElasticity(ctx, elasticity); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}
