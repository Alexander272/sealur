package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_construction_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_base_construction_api"
)

type PutgBaseConstructionService struct {
	repo repository.PutgBaseConstruction
}

func NewPutgBaseConstructionService(repo repository.PutgBaseConstruction) *PutgBaseConstructionService {
	return &PutgBaseConstructionService{
		repo: repo,
	}
}

func (s *PutgBaseConstructionService) Get(ctx context.Context, req *putg_base_construction_api.GetPutgBaseConstruction,
) ([]*putg_construction_type_model.PutgConstruction, error) {
	constructions, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get base construction. error: %w", err)
	}
	return constructions, nil
}

func (s *PutgBaseConstructionService) Create(ctx context.Context, c *putg_base_construction_api.CreatePutgBaseConstruction) error {
	if err := s.repo.Create(ctx, c); err != nil {
		return fmt.Errorf("failed to create base construction. error: %w", err)
	}
	return nil
}

func (s *PutgBaseConstructionService) Update(ctx context.Context, c *putg_base_construction_api.UpdatePutgBaseConstruction) error {
	if err := s.repo.Update(ctx, c); err != nil {
		return fmt.Errorf("failed to update base construction. error: %w", err)
	}
	return nil
}

func (s *PutgBaseConstructionService) Delete(ctx context.Context, c *putg_base_construction_api.DeletePutgBaseConstruction) error {
	if err := s.repo.Delete(ctx, c); err != nil {
		return fmt.Errorf("failed to delete base construction. error: %w", err)
	}
	return nil
}
