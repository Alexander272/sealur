package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
)

type FlangeService struct {
	repo      repository.Flange
	materials *MaterialsService
}

func NewFlangeService(repo repository.Flange, materials *MaterialsService) *FlangeService {
	return &FlangeService{
		repo:      repo,
		materials: materials,
	}
}

func (s *FlangeService) Calculation(ctx context.Context) error {
	size, err := s.repo.GetSize(ctx, 400, 1.0)
	if err != nil {
		return fmt.Errorf("failed to get size. error: %w", err)
	}

	logger.Debug(size)

	return nil
}
