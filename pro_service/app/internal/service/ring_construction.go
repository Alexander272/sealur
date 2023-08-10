package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_construction_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_construction_api"
)

type RingConstructionService struct {
	repo repository.RingConstruction
}

func NewRingConstructionService(repo repository.RingConstruction) *RingConstructionService {
	return &RingConstructionService{
		repo: repo,
	}
}

type RingConstruction interface {
	GetAll(context.Context, *ring_construction_api.GetRingConstructions) (*ring_construction_model.RingConstructionMap, error)
	Create(context.Context, *ring_construction_api.CreateRingConstruction) error
	Update(context.Context, *ring_construction_api.UpdateRingConstruction) error
	Delete(context.Context, *ring_construction_api.DeleteRingConstruction) error
}

func (s *RingConstructionService) GetAll(ctx context.Context, req *ring_construction_api.GetRingConstructions,
) (*ring_construction_model.RingConstructionMap, error) {
	constructions, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all constructions. error: %w", err)
	}
	return constructions, nil
}

func (s *RingConstructionService) Create(ctx context.Context, c *ring_construction_api.CreateRingConstruction) error {
	if err := s.repo.Create(ctx, c); err != nil {
		return fmt.Errorf("failed to create ring construction. error: %w", err)
	}
	return nil
}

func (s *RingConstructionService) Update(ctx context.Context, c *ring_construction_api.UpdateRingConstruction) error {
	if err := s.repo.Update(ctx, c); err != nil {
		return fmt.Errorf("failed to update ring construction. error: %w", err)
	}
	return nil
}

func (s *RingConstructionService) Delete(ctx context.Context, c *ring_construction_api.DeleteRingConstruction) error {
	if err := s.repo.Delete(ctx, c); err != nil {
		return fmt.Errorf("failed to delete ring construction. error: %w", err)
	}
	return nil
}
