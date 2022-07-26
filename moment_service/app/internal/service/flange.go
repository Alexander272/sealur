package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type FlangeService struct {
	repo repository.Flange
}

func NewFlangeService(repo repository.Flange) *FlangeService {
	return &FlangeService{repo: repo}
}

func (s *FlangeService) GetFlangeSize(ctx context.Context, req *moment_api.GetFlangeSizeRequest) (models.FlangeSize, error) {
	size, err := s.repo.GetFlangeSize(ctx, req)
	if err != nil {
		return models.FlangeSize{}, fmt.Errorf("failed to get flange size. error: %w", err)
	}
	size.Area = math.Round(size.Area*1000) / 1000

	return size, nil
}

func (s *FlangeService) CreateFlangeSize(ctx context.Context, size *moment_api.CreateFlangeSizeRequest) error {
	if err := s.repo.CreateFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to create flange size. error: %w", err)
	}
	return nil
}

func (s *FlangeService) UpdateFlangeSize(ctx context.Context, size *moment_api.UpdateFlangeSizeRequest) error {
	if err := s.repo.UpdateFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to update flange size. error: %w", err)
	}
	return nil
}

func (s *FlangeService) DeleteFlangeSize(ctx context.Context, size *moment_api.DeleteFlangeSizeRequest) error {
	if err := s.repo.DeleteFlangeSize(ctx, size); err != nil {
		return fmt.Errorf("failed to delete flange size. error: %w", err)
	}
	return nil
}
