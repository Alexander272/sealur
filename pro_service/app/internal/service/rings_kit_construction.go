package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/rings_kit_construction_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_construction_api"
)

type RingsKitConstructionService struct {
	repo repository.RingsKitConstruction
}

func NewRingsKitConstructionService(repo repository.RingsKitConstruction) *RingsKitConstructionService {
	return &RingsKitConstructionService{
		repo: repo,
	}
}

type RingsKitConstruction interface {
	GetAll(context.Context, *rings_kit_construction_api.GetRingsKitConstructions) (*rings_kit_construction_model.RingsKitConstructionMap, error)
	Create(context.Context, *rings_kit_construction_api.CreateRingsKitConstruction) error
	Update(context.Context, *rings_kit_construction_api.UpdateRingsKitConstruction) error
	Delete(context.Context, *rings_kit_construction_api.DeleteRingsKitConstruction) error
}

func (s *RingsKitConstructionService) GetAll(ctx context.Context, req *rings_kit_construction_api.GetRingsKitConstructions,
) (*rings_kit_construction_model.RingsKitConstructionMap, error) {
	constructions, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all rings kit constructions. error: %w", err)
	}
	return constructions, nil
}

func (s *RingsKitConstructionService) Create(ctx context.Context, c *rings_kit_construction_api.CreateRingsKitConstruction) error {
	if err := s.repo.Create(ctx, c); err != nil {
		return fmt.Errorf("failed to create rings kit construction. error: %w", err)
	}
	return nil
}

func (s *RingsKitConstructionService) Update(ctx context.Context, c *rings_kit_construction_api.UpdateRingsKitConstruction) error {
	if err := s.repo.Update(ctx, c); err != nil {
		return fmt.Errorf("failed to update rings kit construction. error: %w", err)
	}
	return nil
}

func (s *RingsKitConstructionService) Delete(ctx context.Context, c *rings_kit_construction_api.DeleteRingsKitConstruction) error {
	if err := s.repo.Delete(ctx, c); err != nil {
		return fmt.Errorf("failed to delete rings kit construction. error: %w", err)
	}
	return nil
}
