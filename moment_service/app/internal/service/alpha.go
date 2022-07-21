package service

import (
	"context"
	"fmt"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (s *MaterialsService) CreateAlpha(ctx context.Context, alpha *moment_proto.CreateAlphaRequest) error {
	err := s.repo.CreateAlpha(ctx, alpha)
	if err != nil {
		return fmt.Errorf("failed to create alpha. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) UpateAlpha(ctx context.Context, alpha *moment_proto.UpdateAlphaRequest) error {
	if err := s.repo.UpateAlpha(ctx, alpha); err != nil {
		return fmt.Errorf("failed to update alpha. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteAlpha(ctx context.Context, alpha *moment_proto.DeleteAlphaRequest) error {
	if err := s.repo.DeleteAlpha(ctx, alpha); err != nil {
		return fmt.Errorf("failed to delete alpha. error: %w", err)
	}
	return nil
}
