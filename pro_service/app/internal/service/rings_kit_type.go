package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/rings_kit_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_type_api"
)

type RingsKitTypeService struct {
	repo repository.RingsKitType
}

func NewRingsKitTypeService(repo repository.RingsKitType) *RingsKitTypeService {
	return &RingsKitTypeService{
		repo: repo,
	}
}

type RingsKitType interface {
	GetAll(context.Context, *rings_kit_type_api.GetRingsKitTypes) ([]*rings_kit_type_model.RingsKitType, error)
	Create(context.Context, *rings_kit_type_api.CreateRingsKitType) error
	Update(context.Context, *rings_kit_type_api.UpdateRingsKitType) error
	Delete(context.Context, *rings_kit_type_api.DeleteRingsKitType) error
}

func (s *RingsKitTypeService) GetAll(ctx context.Context, req *rings_kit_type_api.GetRingsKitTypes) ([]*rings_kit_type_model.RingsKitType, error) {
	kit, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all rings kit types. error: %w", err)
	}
	return kit, nil
}

func (s *RingsKitTypeService) Create(ctx context.Context, kit *rings_kit_type_api.CreateRingsKitType) error {
	if err := s.repo.Create(ctx, kit); err != nil {
		return fmt.Errorf("failed to create rings kit type. error: %w", err)
	}
	return nil
}

func (s *RingsKitTypeService) Update(ctx context.Context, kit *rings_kit_type_api.UpdateRingsKitType) error {
	if err := s.repo.Update(ctx, kit); err != nil {
		return fmt.Errorf("failed to update rings kit type. error: %w", err)
	}
	return nil
}

func (s *RingsKitTypeService) Delete(ctx context.Context, kit *rings_kit_type_api.DeleteRingsKitType) error {
	if err := s.repo.Delete(ctx, kit); err != nil {
		return fmt.Errorf("failed to delete rings kit type. error: %w", err)
	}
	return nil
}
