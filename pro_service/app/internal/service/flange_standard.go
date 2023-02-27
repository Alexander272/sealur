package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/flange_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_standard_model"
)

type FlangeStandardService struct {
	repo repository.FlangeStandard
}

func NewFlangeStandardService(repo repository.FlangeStandard) *FlangeStandardService {
	return &FlangeStandardService{repo: repo}
}

func (s *FlangeStandardService) GetAll(ctx context.Context, flange *flange_standard_api.GetAllFlangeStandards) ([]*flange_standard_model.FlangeStandard, error) {
	flanges, err := s.repo.GetAll(ctx, flange)
	if err != nil {
		return nil, fmt.Errorf("failed to get flange standard. error: %w", err)
	}
	return flanges, err
}

func (s *FlangeStandardService) Create(ctx context.Context, flange *flange_standard_api.CreateFlangeStandard) error {
	if err := s.repo.Create(ctx, flange); err != nil {
		return fmt.Errorf("failed to create flange standard. error: %w", err)
	}
	return nil
}

func (s *FlangeStandardService) CreateSeveral(ctx context.Context, flanges *flange_standard_api.CreateSeveralFlangeStandard) error {
	if err := s.repo.CreateSeveral(ctx, flanges); err != nil {
		return fmt.Errorf("failed to create several flange standards. error: %w", err)
	}
	return nil
}

func (s *FlangeStandardService) Update(ctx context.Context, flange *flange_standard_api.UpdateFlangeStandard) error {
	if err := s.repo.Update(ctx, flange); err != nil {
		return fmt.Errorf("failed to update flange standard. error: %w", err)
	}
	return nil
}

func (s *FlangeStandardService) Delete(ctx context.Context, flange *flange_standard_api.DeleteFlangeStandard) error {
	if err := s.repo.Delete(ctx, flange); err != nil {
		return fmt.Errorf("failed to delete flange standard. error: %w", err)
	}
	return nil
}
