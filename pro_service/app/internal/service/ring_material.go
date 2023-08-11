package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_material_api"
)

type RingMaterialService struct {
	repo repository.RingMaterial
}

func NewRingMaterialService(repo repository.RingMaterial) *RingMaterialService {
	return &RingMaterialService{
		repo: repo,
	}
}

type RingMaterial interface {
	Get(context.Context, *ring_material_api.GetRingMaterial) ([]*ring_material_model.RingMaterial, error)
	Create(context.Context, *ring_material_api.CreateRingMaterial) error
	Update(context.Context, *ring_material_api.UpdateRingMaterial) error
	Delete(context.Context, *ring_material_api.DeleteRingMaterial) error
}

func (s *RingMaterialService) Get(ctx context.Context, req *ring_material_api.GetRingMaterial) ([]*ring_material_model.RingMaterial, error) {
	materials, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get ring materials. error: %w", err)
	}
	return materials, nil
}

func (s *RingMaterialService) Create(ctx context.Context, m *ring_material_api.CreateRingMaterial) error {
	if err := s.repo.Create(ctx, m); err != nil {
		return fmt.Errorf("failed to create ring material. error: %w", err)
	}
	return nil
}

func (s *RingMaterialService) Update(ctx context.Context, m *ring_material_api.UpdateRingMaterial) error {
	if err := s.repo.Update(ctx, m); err != nil {
		return fmt.Errorf("failed to update ring material. error: %w", err)
	}
	return nil
}

func (s *RingMaterialService) Delete(ctx context.Context, m *ring_material_api.DeleteRingMaterial) error {
	if err := s.repo.Delete(ctx, m); err != nil {
		return fmt.Errorf("failed to delete ring material. error: %w", err)
	}
	return nil
}
