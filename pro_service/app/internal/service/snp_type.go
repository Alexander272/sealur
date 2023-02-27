package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_type_api"
)

type SnpTypeService struct {
	repo repository.SnpType
}

func NewSnpTypeService(repo repository.SnpType) *SnpTypeService {
	return &SnpTypeService{repo: repo}
}

func (s *SnpTypeService) Get(ctx context.Context, flange *snp_type_api.GetSnpTypes) ([]*snp_type_model.SnpType, error) {
	types, err := s.repo.Get(ctx, flange)
	if err != nil {
		return nil, fmt.Errorf("failed to get snp type. error: %w", err)
	}
	return types, err
}

func (s *SnpTypeService) Create(ctx context.Context, flange *snp_type_api.CreateSnpType) error {
	if err := s.repo.Create(ctx, flange); err != nil {
		return fmt.Errorf("failed to create snp type. error: %w", err)
	}
	return nil
}

func (s *SnpTypeService) CreateSeveral(ctx context.Context, types *snp_type_api.CreateSeveralSnpType) error {
	if err := s.repo.CreateSeveral(ctx, types); err != nil {
		return fmt.Errorf("failed to create several snp types. error: %w", err)
	}
	return nil
}

func (s *SnpTypeService) Update(ctx context.Context, flange *snp_type_api.UpdateSnpType) error {
	if err := s.repo.Update(ctx, flange); err != nil {
		return fmt.Errorf("failed to update snp type. error: %w", err)
	}
	return nil
}

func (s *SnpTypeService) Delete(ctx context.Context, flange *snp_type_api.DeleteSnpType) error {
	if err := s.repo.Delete(ctx, flange); err != nil {
		return fmt.Errorf("failed to delete snp type. error: %w", err)
	}
	return nil
}
