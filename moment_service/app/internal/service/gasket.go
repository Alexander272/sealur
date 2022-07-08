package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
)

type GasketService struct {
	repo repository.Gasket
}

func NewGasketService(repo repository.Gasket) *GasketService {
	return &GasketService{
		repo: repo,
	}
}

func (s *GasketService) Get(ctx context.Context, gasket models.GetGasket) (models.Gasket, error) {
	g, err := s.repo.Get(ctx, gasket)
	if err != nil {
		return models.Gasket{}, fmt.Errorf("failed to get gasket. error: %w", err)
	}
	return g, nil
}
