package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/snp_data_api"
)

type SnpDataService struct {
	repo repository.SnpData
}

func NewSnpDataService(repo repository.SnpData) *SnpDataService {
	return &SnpDataService{repo: repo}
}

// func (s *SnpDataService) Get(ctx context.Context, req *snp_data_api.GetSnpData) ([]*snp_data_model.SnpData, error) {
// 	snp, err := s.repo.Get(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get snp standard. error: %w", err)
// 	}
// 	return snp, err
// }

func (s *SnpDataService) Create(ctx context.Context, snp *snp_data_api.CreateSnpData) error {
	if err := s.repo.Create(ctx, snp); err != nil {
		return fmt.Errorf("failed to create snp standard. error: %w", err)
	}
	return nil
}

func (s *SnpDataService) Update(ctx context.Context, snp *snp_data_api.UpdateSnpData) error {
	if err := s.repo.Update(ctx, snp); err != nil {
		return fmt.Errorf("failed to update snp standard. error: %w", err)
	}
	return nil
}

func (s *SnpDataService) Delete(ctx context.Context, snp *snp_data_api.DeleteSnpData) error {
	if err := s.repo.Delete(ctx, snp); err != nil {
		return fmt.Errorf("failed to delete snp standard. error: %w", err)
	}
	return nil
}
