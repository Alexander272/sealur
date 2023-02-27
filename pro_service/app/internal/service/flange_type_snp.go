package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_type_snp_model"
)

type FlangeTypeSnpService struct {
	repo repository.FlangeTypeSnp
}

func NewFlangeTypeSnpService(repo repository.FlangeTypeSnp) *FlangeTypeSnpService {
	return &FlangeTypeSnpService{repo: repo}
}

func (s *FlangeTypeSnpService) Get(ctx context.Context, flange *flange_type_snp_api.GetFlangeTypeSnp) ([]*flange_type_snp_model.FlangeTypeSnp, error) {
	flanges, err := s.repo.Get(ctx, flange)
	if err != nil {
		return nil, fmt.Errorf("failed to get flange type snp. error: %w", err)
	}
	return flanges, err
}

func (s *FlangeTypeSnpService) Create(ctx context.Context, flange *flange_type_snp_api.CreateFlangeTypeSnp) error {
	if err := s.repo.Create(ctx, flange); err != nil {
		return fmt.Errorf("failed to create flange type snp. error: %w", err)
	}
	return nil
}

func (s *FlangeTypeSnpService) CreateSeveral(ctx context.Context, flanges *flange_type_snp_api.CreateSeveralFlangeTypeSnp) error {
	if err := s.repo.CreateSeveral(ctx, flanges); err != nil {
		return fmt.Errorf("failed to create several flange type snp. error: %w", err)
	}
	return nil
}

func (s *FlangeTypeSnpService) Update(ctx context.Context, flange *flange_type_snp_api.UpdateFlangeTypeSnp) error {
	if err := s.repo.Update(ctx, flange); err != nil {
		return fmt.Errorf("failed to update flange type snp. error: %w", err)
	}
	return nil
}

func (s *FlangeTypeSnpService) Delete(ctx context.Context, flange *flange_type_snp_api.DeleteFlangeTypeSnp) error {
	if err := s.repo.Delete(ctx, flange); err != nil {
		return fmt.Errorf("failed to delete flange type snp. error: %w", err)
	}
	return nil
}
