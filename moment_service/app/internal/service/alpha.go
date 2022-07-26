package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (s *MaterialsService) CreateAlpha(ctx context.Context, alpha *moment_api.CreateAlphaRequest) error {
	err := s.repo.CreateAlpha(ctx, alpha)
	if err != nil {
		return fmt.Errorf("failed to create alpha. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) UpateAlpha(ctx context.Context, alpha *moment_api.UpdateAlphaRequest) error {
	if err := s.repo.UpateAlpha(ctx, alpha); err != nil {
		return fmt.Errorf("failed to update alpha. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteAlpha(ctx context.Context, alpha *moment_api.DeleteAlphaRequest) error {
	if err := s.repo.DeleteAlpha(ctx, alpha); err != nil {
		return fmt.Errorf("failed to delete alpha. error: %w", err)
	}
	return nil
}
