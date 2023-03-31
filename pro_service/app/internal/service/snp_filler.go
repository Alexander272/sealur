package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
)

type SnpFillerService struct {
	repo repository.SnpFiller
}

func NewSnpFillerService(repo repository.SnpFiller) *SnpFillerService {
	return &SnpFillerService{repo: repo}
}

func (s *SnpFillerService) GetAll(ctx context.Context, req *snp_filler_api.GetSnpFillers) ([]*snp_filler_model.SnpFiller, error) {
	fillers, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get snp filler. error: %w", err)
	}
	return fillers, err
}

func (s *SnpFillerService) GetAllNew(ctx context.Context, req *snp_filler_api.GetSnpFillers) ([]*snp_filler_model.SnpFillerNew, error) {
	fillers, err := s.repo.GetAllNew(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get snp filler new. error: %w", err)
	}
	return fillers, err
}

func (s *SnpFillerService) Create(ctx context.Context, filler *snp_filler_api.CreateSnpFiller) error {
	if err := s.repo.Create(ctx, filler); err != nil {
		return fmt.Errorf("failed to create snp filler. error: %w", err)
	}
	return nil
}

func (s *SnpFillerService) CreateSeveral(ctx context.Context, fillers *snp_filler_api.CreateSeveralSnpFiller) error {
	if err := s.repo.CreateSeveral(ctx, fillers); err != nil {
		return fmt.Errorf("failed to create several snp fillers. error: %w", err)
	}
	return nil
}

func (s *SnpFillerService) Update(ctx context.Context, filler *snp_filler_api.UpdateSnpFiller) error {
	if err := s.repo.Update(ctx, filler); err != nil {
		return fmt.Errorf("failed to update snp filler. error: %w", err)
	}
	return nil
}

func (s *SnpFillerService) Delete(ctx context.Context, filler *snp_filler_api.DeleteSnpFiller) error {
	if err := s.repo.Delete(ctx, filler); err != nil {
		return fmt.Errorf("failed to delete snp filler. error: %w", err)
	}
	return nil
}
