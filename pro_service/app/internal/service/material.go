package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/material_model"
)

type MaterialService struct {
	repo repository.Material
}

func NewMaterialService(repo repository.Material) *MaterialService {
	return &MaterialService{repo: repo}
}

func (s *MaterialService) GetAll(ctx context.Context, material *material_api.GetAllMaterials) ([]*material_model.Material, error) {
	materials, err := s.repo.GetAll(ctx, material)
	if err != nil {
		return nil, fmt.Errorf("failed to get material. error: %w", err)
	}
	return materials, err
}

func (s *MaterialService) Create(ctx context.Context, material *material_api.CreateMaterial) error {
	if err := s.repo.Create(ctx, material); err != nil {
		return fmt.Errorf("failed to create material. error: %w", err)
	}
	return nil
}

func (s *MaterialService) CreateSeveral(ctx context.Context, materials *material_api.CreateSeveralMaterial) error {
	if err := s.repo.CreateSeveral(ctx, materials); err != nil {
		return fmt.Errorf("failed to create several materials. error: %w", err)
	}
	return nil
}

func (s *MaterialService) Update(ctx context.Context, material *material_api.UpdateMaterial) error {
	if err := s.repo.Update(ctx, material); err != nil {
		return fmt.Errorf("failed to update material. error: %w", err)
	}
	return nil
}

func (s *MaterialService) Delete(ctx context.Context, material *material_api.DeleteMaterial) error {
	if err := s.repo.Delete(ctx, material); err != nil {
		return fmt.Errorf("failed to delete material. error: %w", err)
	}
	return nil
}
