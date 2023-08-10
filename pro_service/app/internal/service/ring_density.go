package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_density_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_density_api"
)

type RingDensityService struct {
	repo repository.RingDensity
}

func NewRingDensityService(repo repository.RingDensity) *RingDensityService {
	return &RingDensityService{
		repo: repo,
	}
}

type RingDensity interface {
	GetAll(context.Context, *ring_density_api.GetRingDensity) (*ring_density_model.RingDensityMap, error)
	Create(context.Context, *ring_density_api.CreateRingDensity) error
	Update(context.Context, *ring_density_api.UpdateRingDensity) error
	Delete(context.Context, *ring_density_api.DeleteRingDensity) error
}

func (s *RingDensityService) GetAll(ctx context.Context, req *ring_density_api.GetRingDensity) (*ring_density_model.RingDensityMap, error) {
	ringTypes, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ring density. error: %w", err)
	}
	return ringTypes, nil
}

func (s *RingDensityService) Create(ctx context.Context, density *ring_density_api.CreateRingDensity) error {
	if err := s.repo.Create(ctx, density); err != nil {
		return fmt.Errorf("failed to create ring density. error: %w", err)
	}
	return nil
}

func (s *RingDensityService) Update(ctx context.Context, density *ring_density_api.UpdateRingDensity) error {
	if err := s.repo.Update(ctx, density); err != nil {
		return fmt.Errorf("failed to update ring density. error: %w", err)
	}
	return nil
}

func (s *RingDensityService) Delete(ctx context.Context, density *ring_density_api.DeleteRingDensity) error {
	if err := s.repo.Delete(ctx, density); err != nil {
		return fmt.Errorf("failed to delete ring density. error: %w", err)
	}
	return nil
}
