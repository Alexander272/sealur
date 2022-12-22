package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/material_api"
)

func (h *MaterialsHandlers) CreateElasticity(ctx context.Context, elasticity *material_api.CreateElasticityRequest) (*material_api.Response, error) {
	err := h.service.CreateElasticity(ctx, elasticity)
	if err != nil {
		return nil, err
	}

	return &material_api.Response{}, nil
}

func (h *MaterialsHandlers) UpdateElasticity(ctx context.Context, elasticity *material_api.UpdateElasticityRequest) (*material_api.Response, error) {
	if err := h.service.UpdateElasticity(ctx, elasticity); err != nil {
		return nil, err
	}
	return &material_api.Response{}, nil
}

func (h *MaterialsHandlers) DeleteElasticity(ctx context.Context, elasticity *material_api.DeleteElasticityRequest) (*material_api.Response, error) {
	if err := h.service.DeleteElasticity(ctx, elasticity); err != nil {
		return nil, err
	}
	return &material_api.Response{}, nil
}
