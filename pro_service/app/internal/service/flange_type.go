package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_type_model"
)

type FlangeTypeService struct {
	repo repository.FlangeType
}

func NewFlangeTypeService(repo repository.FlangeType) *FlangeTypeService {
	return &FlangeTypeService{
		repo: repo,
	}
}

func (s *FlangeTypeService) Get(ctx context.Context, req *flange_type_api.GetFlangeType) ([]*flange_type_model.FlangeType, error) {
	flangeTypes, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get flange type. error: %w", err)
	}
	return flangeTypes, nil
}

func (s *FlangeTypeService) Create(ctx context.Context, flangeType *flange_type_api.CreateFlangeType) error {
	if err := s.repo.Create(ctx, flangeType); err != nil {
		return fmt.Errorf("failed to create flange type. error: %w", err)
	}
	return nil
}

func (s *FlangeTypeService) Update(ctx context.Context, flangeType *flange_type_api.UpdateFlangeType) error {
	if err := s.repo.Update(ctx, flangeType); err != nil {
		return fmt.Errorf("failed to update flange type. error: %w", err)
	}
	return nil
}

func (s *FlangeTypeService) Delete(ctx context.Context, flangeType *flange_type_api.DeleteFlangeType) error {
	if err := s.repo.Delete(ctx, flangeType); err != nil {
		return fmt.Errorf("failed to delete flange type. error: %w", err)
	}
	return nil
}
