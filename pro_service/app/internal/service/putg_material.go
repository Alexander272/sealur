package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_material_api"
)

type PutgMaterialService struct {
	repo repository.PutgMaterial
}

func NewPutgMaterialService(repo repository.PutgMaterial) *PutgMaterialService {
	return &PutgMaterialService{
		repo: repo,
	}
}

func (s *PutgMaterialService) Get(ctx context.Context, req *putg_material_api.GetPutgMaterial) (*putg_material_model.PutgMaterials, error) {
	materials, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg materials. error: %w", err)
	}
	return materials, nil
}

//TODO добавить создание материалов (которые в MaterialId идут)

func (s *PutgMaterialService) Create(ctx context.Context, material *putg_material_api.CreatePutgMaterial) error {
	if err := s.repo.Create(ctx, material); err != nil {
		return fmt.Errorf("failed to create putg material. error: %w", err)
	}
	return nil
}

func (s *PutgMaterialService) Update(ctx context.Context, material *putg_material_api.UpdatePutgMaterial) error {
	if err := s.repo.Update(ctx, material); err != nil {
		return fmt.Errorf("failed to update putg material. error: %w", err)
	}
	return nil
}

func (s *PutgMaterialService) Delete(ctx context.Context, material *putg_material_api.DeletePutgMaterial) error {
	if err := s.repo.Delete(ctx, material); err != nil {
		return fmt.Errorf("failed to delete putg material. error: %w", err)
	}
	return nil
}
