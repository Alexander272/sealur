package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_size_api"
)

type SnpSizeService struct {
	repo repository.SnpSize
}

func NewSnpSizeService(repo repository.SnpSize) *SnpSizeService {
	return &SnpSizeService{
		repo: repo,
	}
}

func (s *SnpSizeService) Get(ctx context.Context, req *snp_size_api.GetSnpSize) ([]*snp_size_model.SnpSize, error) {
	sizes, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get snp sizes. error: %w", err)
	}
	return sizes, nil
}

func (s *SnpSizeService) Create(ctx context.Context, size *snp_size_api.CreateSnpSize) error {
	if err := s.repo.Create(ctx, size); err != nil {
		return fmt.Errorf("failed to create snp size. error: %w", err)
	}
	return nil
}
