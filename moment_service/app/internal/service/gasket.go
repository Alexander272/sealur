package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type GasketService struct {
	repo repository.Gasket
}

func NewGasketService(repo repository.Gasket) *GasketService {
	return &GasketService{
		repo: repo,
	}
}

func (s *GasketService) GetFullData(ctx context.Context, gasket models.GetGasket) (models.FullDataGasket, error) {
	g, err := s.repo.GetFullData(ctx, gasket)
	if err != nil {
		return models.FullDataGasket{}, fmt.Errorf("failed to get gasket. error: %w", err)
	}
	return g, nil
}

func (s *GasketService) GetGasket(ctx context.Context, req *moment_api.GetGasketRequest) (gasket []*moment_api.Gasket, err error) {
	data, err := s.repo.GetGasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get gasket. error: %w", err)
	}

	for _, item := range data {
		gasket = append(gasket, &moment_api.Gasket{
			Id:    item.Id,
			Title: item.Title,
		})
	}

	return gasket, nil
}

func (s *GasketService) GetGasketWithThick(ctx context.Context, req *moment_api.GetGasketRequest) (gasket []*moment_api.GasketWithThick, err error) {
	data, err := s.repo.GetGasketWithThick(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get gasket. error: %w", err)
	}

	curId := ""
	for _, item := range data {
		item.Thickness = math.Round(item.Thickness*1000) / 1000
		if item.Id != curId {
			curId = item.Id
			gasket = append(gasket, &moment_api.GasketWithThick{
				Id:    item.Id,
				Title: item.Title,
			})
			gasket[len(gasket)-1].Thickness = append(gasket[len(gasket)-1].Thickness, item.Thickness)
		} else {
			gasket[len(gasket)-1].Thickness = append(gasket[len(gasket)-1].Thickness, item.Thickness)
		}
	}

	return gasket, nil
}

func (s *GasketService) CreateGasket(ctx context.Context, gasket *moment_api.CreateGasketRequest) (id string, err error) {
	id, err = s.repo.CreateGasket(ctx, gasket)
	if err != nil {
		return "", fmt.Errorf("failed to create gasket. error: %w", err)
	}
	return id, nil
}

func (s *GasketService) UpdateGasket(ctx context.Context, gasket *moment_api.UpdateGasketRequest) error {
	if err := s.repo.UpdateGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to update gasket. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteGasket(ctx context.Context, gasket *moment_api.DeleteGasketRequest) error {
	if err := s.repo.DeleteGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to delete gasket. error: %w", err)
	}
	return nil
}

//---

func (s *GasketService) CreateGasketData(ctx context.Context, data *moment_api.CreateGasketDataRequest) error {
	if err := s.repo.CreateGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to create gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) UpdateGasketData(ctx context.Context, data *moment_api.UpdateGasketDataRequest) error {
	if err := s.repo.UpdateGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to update gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteGasketData(ctx context.Context, data *moment_api.DeleteGasketDataRequest) error {
	if err := s.repo.DeleteGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to delete gasket data. error: %w", err)
	}
	return nil
}
