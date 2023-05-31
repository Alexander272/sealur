package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_size_api"
)

type PutgSizeService struct {
	repo repository.PutgSize
}

func NewPutgSizeService(repo repository.PutgSize) *PutgSizeService {
	return &PutgSizeService{
		repo: repo,
	}
}

func (s *PutgSizeService) Get(ctx context.Context, req *putg_size_api.GetPutgSize) ([]*putg_size_model.PutgSize, error) {
	sizes, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg sizes. error: %w", err)
	}
	return sizes, nil
}

func (s *PutgSizeService) Create(ctx context.Context, size *putg_size_api.CreatePutgSize) error {
	if err := s.repo.Create(ctx, size); err != nil {
		return fmt.Errorf("failed to create putg size. error: %w", err)
	}
	return nil
}

func (s *PutgSizeService) CreateSeveral(ctx context.Context, sizes *putg_size_api.CreateSeveralPutgSize) error {
	if err := s.repo.CreateSeveral(ctx, sizes); err != nil {
		return fmt.Errorf("failed to create several putg sizes. error: %w", err)
	}
	return nil
}

func (s *PutgSizeService) Update(ctx context.Context, size *putg_size_api.UpdatePutgSize) error {
	if err := s.repo.Update(ctx, size); err != nil {
		return fmt.Errorf("failed to update putg size. error: %w", err)
	}
	return nil
}

func (s *PutgSizeService) Delete(ctx context.Context, size *putg_size_api.DeletePutgSize) error {
	if err := s.repo.Delete(ctx, size); err != nil {
		return fmt.Errorf("failed to delete putg size. error: %w", err)
	}
	return nil
}
