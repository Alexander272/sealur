package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/temperature_model"
	"github.com/Alexander272/sealur_proto/api/pro/temperature_api"
)

type TemperatureService struct {
	repo repository.Temperature
}

func NewTemperatureService(repo repository.Temperature) *TemperatureService {
	return &TemperatureService{repo: repo}
}

func (s *TemperatureService) GetAll(ctx context.Context, temperature *temperature_api.GetAllTemperatures) ([]*temperature_model.Temperature, error) {
	temperatures, err := s.repo.GetAll(ctx, temperature)
	if err != nil {
		return nil, fmt.Errorf("failed to get temperature. error: %w", err)
	}
	return temperatures, err
}

func (s *TemperatureService) Create(ctx context.Context, temperature *temperature_api.CreateTemperature) error {
	if err := s.repo.Create(ctx, temperature); err != nil {
		return fmt.Errorf("failed to create temperature. error: %w", err)
	}
	return nil
}

func (s *TemperatureService) CreateSeveral(ctx context.Context, temperatures *temperature_api.CreateSeveralTemperature) error {
	if err := s.repo.CreateSeveral(ctx, temperatures); err != nil {
		return fmt.Errorf("failed to create several temperatures. error: %w", err)
	}
	return nil
}

func (s *TemperatureService) Update(ctx context.Context, temperature *temperature_api.UpdateTemperature) error {
	if err := s.repo.Update(ctx, temperature); err != nil {
		return fmt.Errorf("failed to update temperature. error: %w", err)
	}
	return nil
}

func (s *TemperatureService) Delete(ctx context.Context, temperature *temperature_api.DeleteTemperature) error {
	if err := s.repo.Delete(ctx, temperature); err != nil {
		return fmt.Errorf("failed to delete temperature. error: %w", err)
	}
	return nil
}
