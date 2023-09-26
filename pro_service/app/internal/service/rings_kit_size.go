package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/rings_kit_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_size_api"
)

type RingsKitSizeService struct {
	repo repository.RingsKitSize
}

func NewRingsKitSizeService(repo repository.RingsKitSize) *RingsKitSizeService {
	return &RingsKitSizeService{
		repo: repo,
	}
}

type RingsKitSize interface {
	Get(context.Context, *rings_kit_size_api.GetRingsKitSize) ([]*rings_kit_size_model.RingsKitSize, error)
	Create(context.Context, *rings_kit_size_api.CreateRingsKitSize) error
	Update(context.Context, *rings_kit_size_api.UpdateRingsKitSize) error
	Delete(context.Context, *rings_kit_size_api.DeleteRingsKitSize) error
}

func (s *RingsKitSizeService) Get(ctx context.Context, req *rings_kit_size_api.GetRingsKitSize) ([]*rings_kit_size_model.RingsKitSize, error) {
	sizes, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get rings kit sizes by construction id. error: %w", err)
	}
	return sizes, nil
}

func (s *RingsKitSizeService) Create(ctx context.Context, size *rings_kit_size_api.CreateRingsKitSize) error {
	if err := s.repo.Create(ctx, size); err != nil {
		return fmt.Errorf("failed to create rings kit size. error: %w", err)
	}
	return nil
}

func (s *RingsKitSizeService) Update(ctx context.Context, size *rings_kit_size_api.UpdateRingsKitSize) error {
	if err := s.repo.Update(ctx, size); err != nil {
		return fmt.Errorf("failed to update rings kit size. error: %w", err)
	}
	return nil
}

func (s *RingsKitSizeService) Delete(ctx context.Context, size *rings_kit_size_api.DeleteRingsKitSize) error {
	if err := s.repo.Delete(ctx, size); err != nil {
		return fmt.Errorf("failed to delete rings kit size. error: %w", err)
	}
	return nil
}
