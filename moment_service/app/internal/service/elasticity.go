package service

import (
	"context"
	"fmt"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (s *MaterialsService) CreateElasticity(ctx context.Context, elasticity *moment_proto.CreateElasticityRequest) error {
	err := s.repo.CreateElasticity(ctx, elasticity)
	if err != nil {
		return fmt.Errorf("failed to create elasticity. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) UpdateElasticity(ctx context.Context, elasticity *moment_proto.UpdateElasticityRequest) error {
	if err := s.repo.UpdateElasticity(ctx, elasticity); err != nil {
		return fmt.Errorf("failed to update elasticity. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteElasticity(ctx context.Context, elasticity *moment_proto.DeleteElasticityRequest) error {
	if err := s.repo.DeleteElasticity(ctx, elasticity); err != nil {
		return fmt.Errorf("failed to delete elasticity. error: %w", err)
	}
	return nil
}
