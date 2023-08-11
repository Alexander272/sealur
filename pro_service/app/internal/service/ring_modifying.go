package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_modifying_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_modifying_api"
)

type RingModifyingService struct {
	repo repository.RingModifying
}

func NewRingModifyingService(repo repository.RingModifying) *RingModifyingService {
	return &RingModifyingService{
		repo: repo,
	}
}

type RingModifying interface {
	GetAll(context.Context, *ring_modifying_api.GetRingModifying) ([]*ring_modifying_model.RingModifying, error)
	Create(context.Context, *ring_modifying_api.CreateRingModifying) error
	Update(context.Context, *ring_modifying_api.UpdateRingModifying) error
	Delete(context.Context, *ring_modifying_api.DeleteRingModifying) error
}

func (s *RingModifyingService) GetAll(ctx context.Context, req *ring_modifying_api.GetRingModifying) ([]*ring_modifying_model.RingModifying, error) {
	mods, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ring modifying. error: %w", err)
	}
	return mods, nil
}

func (s *RingModifyingService) Create(ctx context.Context, m *ring_modifying_api.CreateRingModifying) error {
	if err := s.repo.Create(ctx, m); err != nil {
		return fmt.Errorf("failed to create ring modifying. error: %w", err)
	}
	return nil
}

func (s *RingModifyingService) Update(ctx context.Context, m *ring_modifying_api.UpdateRingModifying) error {
	if err := s.repo.Update(ctx, m); err != nil {
		return fmt.Errorf("failed to update ring modifying. error: %w", err)
	}
	return nil
}

func (s *RingModifyingService) Delete(ctx context.Context, m *ring_modifying_api.DeleteRingModifying) error {
	if err := s.repo.Delete(ctx, m); err != nil {
		return fmt.Errorf("failed to delete ring modifying. error: %w", err)
	}
	return nil
}
