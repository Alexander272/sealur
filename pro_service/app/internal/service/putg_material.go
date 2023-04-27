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
