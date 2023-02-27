package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
)

type SnpMaterialService struct {
	repo repository.SnpMaterial
}

func NewSnpMaterialService(repo repository.SnpMaterial) *SnpMaterialService {
	return &SnpMaterialService{repo: repo}
}

// func (s *SnpMaterialService) Get(ctx context.Context, material *snp_material_api.GetAllSnpMaterials) ([]*material_standard_model.SnpMaterial, error) {
// 	materials, err := s.repo.GetAll(ctx, material)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get snp material. error: %w", err)
// 	}
// 	return materials, err
// }

func (s *SnpMaterialService) Create(ctx context.Context, material *snp_material_api.CreateSnpMaterial) error {
	if err := s.repo.Create(ctx, material); err != nil {
		return fmt.Errorf("failed to create snp material. error: %w", err)
	}
	return nil
}

func (s *SnpMaterialService) Update(ctx context.Context, material *snp_material_api.UpdateSnpMaterial) error {
	if err := s.repo.Update(ctx, material); err != nil {
		return fmt.Errorf("failed to update snp material. error: %w", err)
	}
	return nil
}

func (s *SnpMaterialService) Delete(ctx context.Context, material *snp_material_api.DeleteSnpMaterial) error {
	if err := s.repo.Delete(ctx, material); err != nil {
		return fmt.Errorf("failed to delete snp material. error: %w", err)
	}
	return nil
}
