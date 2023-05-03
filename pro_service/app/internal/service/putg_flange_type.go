package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_flange_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_flange_type_api"
)

type PutgFlangeTypeService struct {
	repo repository.PutgFlangeType
}

func NewPutgFlangeTypeService(repo repository.PutgFlangeType) *PutgFlangeTypeService {
	return &PutgFlangeTypeService{
		repo: repo,
	}
}

func (s *PutgFlangeTypeService) Get(ctx context.Context, req *putg_flange_type_api.GetPutgFlangeType) ([]*putg_flange_type_model.PutgFlangeType, error) {
	flangeTypes, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg flange types. error: %w", err)
	}
	return flangeTypes, nil
}

func (s *PutgFlangeTypeService) Create(ctx context.Context, flangeType *putg_flange_type_api.CreatePutgFlangeType) error {
	if err := s.repo.Create(ctx, flangeType); err != nil {
		return fmt.Errorf("failed to create putg flange type. error: %w", err)
	}
	return nil
}

func (s *PutgFlangeTypeService) Update(ctx context.Context, flangeType *putg_flange_type_api.UpdatePutgFlangeType) error {
	if err := s.repo.Update(ctx, flangeType); err != nil {
		return fmt.Errorf("failed to update putg flange type. error: %w", err)
	}
	return nil
}

func (s *PutgFlangeTypeService) Delete(ctx context.Context, flangeType *putg_flange_type_api.DeletePutgFlangeType) error {
	if err := s.repo.Delete(ctx, flangeType); err != nil {
		return fmt.Errorf("failed to delete putg flange type. error: %w", err)
	}
	return nil
}
