package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_size_api"
)

type RingSizeService struct {
	repo repository.RingSize
}

func NewRingSizeService(repo repository.RingSize) *RingSizeService {
	return &RingSizeService{
		repo: repo,
	}
}

type RingSize interface {
	GetAll(context.Context, *ring_size_api.GetRingSize) ([]*ring_size_model.RingSize, error)
	Create(context.Context, *ring_size_api.CreateRingSize) error
	Update(context.Context, *ring_size_api.UpdateRingSize) error
	Delete(context.Context, *ring_size_api.DeleteRingSize) error
}

func (s *RingSizeService) GetAll(ctx context.Context, req *ring_size_api.GetRingSize) ([]*ring_size_model.RingSize, error) {
	sizes, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ring size. error: %w", err)
	}
	return sizes, nil
}

func (s *RingSizeService) Create(ctx context.Context, size *ring_size_api.CreateRingSize) error {
	if err := s.repo.Create(ctx, size); err != nil {
		return fmt.Errorf("failed to create ring size. error: %w", err)
	}
	return nil
}

func (s *RingSizeService) Update(ctx context.Context, size *ring_size_api.UpdateRingSize) error {
	if err := s.repo.Update(ctx, size); err != nil {
		return fmt.Errorf("failed to update ring size. error: %w", err)
	}
	return nil
}

func (s *RingSizeService) Delete(ctx context.Context, size *ring_size_api.DeleteRingSize) error {
	if err := s.repo.Delete(ctx, size); err != nil {
		return fmt.Errorf("failed to delete ring size. error: %w", err)
	}
	return nil
}
