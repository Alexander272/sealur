package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
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

func (s *GasketService) GetGasket(ctx context.Context, req *moment_proto.GetGasketRequest) (gasket []*moment_proto.Gasket, err error) {
	data, err := s.repo.GetGasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get gasket. error: %w", err)
	}

	for _, item := range data {
		g := moment_proto.Gasket(item)
		gasket = append(gasket, &g)
	}

	return gasket, nil
}

func (s *GasketService) CreateGasket(ctx context.Context, gasket *moment_proto.CreateGasketRequest) (id string, err error) {
	id, err = s.repo.CreateGasket(ctx, gasket)
	if err != nil {
		return "", fmt.Errorf("failed to create gasket. error: %w", err)
	}
	return id, nil
}

func (s *GasketService) UpdateGasket(ctx context.Context, gasket *moment_proto.UpdateGasketRequest) error {
	if err := s.repo.UpdateGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to update gasket. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteGasket(ctx context.Context, gasket *moment_proto.DeleteGasketRequest) error {
	if err := s.repo.DeleteGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to delete gasket. error: %w", err)
	}
	return nil
}

//---

func (s *GasketService) CreateGasketData(ctx context.Context, data *moment_proto.CreateGasketDataRequest) error {
	if err := s.repo.CreateGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to create gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) UpdateGasketData(ctx context.Context, data *moment_proto.UpdateGasketDataRequest) error {
	if err := s.repo.UpdateGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to update gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteGasketData(ctx context.Context, data *moment_proto.DeleteGasketDataRequest) error {
	if err := s.repo.DeleteGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to delete gasket data. error: %w", err)
	}
	return nil
}
