package grpc

import (
	"context"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (h *MaterialsHandlers) CreateElasticity(ctx context.Context, elasticity *moment_proto.CreateElasticityRequest) (*moment_proto.Response, error) {
	err := h.service.CreateElasticity(ctx, elasticity)
	if err != nil {
		return nil, err
	}

	return &moment_proto.Response{}, nil
}

func (h *MaterialsHandlers) UpdateElasticity(ctx context.Context, elasticity *moment_proto.UpdateElasticityRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateElasticity(ctx, elasticity); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *MaterialsHandlers) DeleteElasticity(ctx context.Context, elasticity *moment_proto.DeleteElasticityRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteElasticity(ctx, elasticity); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
