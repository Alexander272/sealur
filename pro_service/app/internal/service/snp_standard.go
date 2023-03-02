package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_standard_api"
)

type SnpStandardService struct {
	repo repository.SnpStandard
}

func NewSnpStandardService(repo repository.SnpStandard) *SnpStandardService {
	return &SnpStandardService{repo: repo}
}

func (s *SnpStandardService) GetAll(ctx context.Context, standard *snp_standard_api.GetAllSnpStandards) ([]*snp_standard_model.SnpStandard, error) {
	standards, err := s.repo.GetAll(ctx, standard)
	if err != nil {
		return nil, fmt.Errorf("failed to get snp standard. error: %w", err)
	}
	return standards, nil
}

func (s *SnpStandardService) GetDefault(ctx context.Context) (*snp_standard_model.SnpStandard, error) {
	standard, err := s.repo.GetDefault(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get default snp standard. error: %w", err)
	}
	return standard, nil
}

func (s *SnpStandardService) Create(ctx context.Context, standard *snp_standard_api.CreateSnpStandard) error {
	if err := s.repo.Create(ctx, standard); err != nil {
		return fmt.Errorf("failed to create snp standard. error: %w", err)
	}
	return nil
}

func (s *SnpStandardService) CreateSeveral(ctx context.Context, standard *snp_standard_api.CreateSeveralSnpStandard) error {
	if err := s.repo.CreateSeveral(ctx, standard); err != nil {
		return fmt.Errorf("failed to create several snp standards. error: %w", err)
	}
	return nil
}

func (s *SnpStandardService) Update(ctx context.Context, standard *snp_standard_api.UpdateSnpStandard) error {
	if err := s.repo.Update(ctx, standard); err != nil {
		return fmt.Errorf("failed to update snp standard. error: %w", err)
	}
	return nil
}

func (s *SnpStandardService) Delete(ctx context.Context, standard *snp_standard_api.DeleteSnpStandard) error {
	if err := s.repo.Delete(ctx, standard); err != nil {
		return fmt.Errorf("failed to delete snp standard. error: %w", err)
	}
	return nil
}
