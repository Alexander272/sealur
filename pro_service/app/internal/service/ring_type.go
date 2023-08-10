package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_type_api"
)

type RingTypeService struct {
	repo repository.RingType
}

func NewRingTypeService(repo repository.RingType) *RingTypeService {
	return &RingTypeService{
		repo: repo,
	}
}

type RingType interface {
	GetAll(context.Context, *ring_type_api.GetRingTypes) ([]*ring_type_model.RingType, error)
	Create(context.Context, *ring_type_api.CreateRingType) error
	Update(context.Context, *ring_type_api.UpdateRingType) error
	Delete(context.Context, *ring_type_api.DeleteRingType) error
}

func (s *RingTypeService) GetAll(ctx context.Context, req *ring_type_api.GetRingTypes) ([]*ring_type_model.RingType, error) {
	ringTypes, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ring types. error: %w", err)
	}

	return ringTypes, nil
}

func (s *RingTypeService) Create(ctx context.Context, ring *ring_type_api.CreateRingType) error {
	if err := s.repo.Create(ctx, ring); err != nil {
		return fmt.Errorf("failed to create ring type. error: %w", err)
	}
	return nil
}

func (s *RingTypeService) Update(ctx context.Context, ring *ring_type_api.UpdateRingType) error {
	if err := s.repo.Update(ctx, ring); err != nil {
		return fmt.Errorf("failed to update ring type. error: %w", err)
	}
	return nil
}

func (s *RingTypeService) Delete(ctx context.Context, ring *ring_type_api.DeleteRingType) error {
	if err := s.repo.Delete(ctx, ring); err != nil {
		return fmt.Errorf("failed to delete ring type. error: %w", err)
	}
	return nil
}
