package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/mounting_model"
	"github.com/Alexander272/sealur_proto/api/pro/mounting_api"
)

type MountingService struct {
	repo repository.Mounting
}

func NewMountingService(repo repository.Mounting) *MountingService {
	return &MountingService{repo: repo}
}

func (s *MountingService) GetAll(ctx context.Context, req *mounting_api.GetAllMountings) ([]*mounting_model.Mounting, error) {
	mounting, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get mounting. error: %w", err)
	}
	return mounting, err
}

func (s *MountingService) Create(ctx context.Context, mounting *mounting_api.CreateMounting) error {
	if err := s.repo.Create(ctx, mounting); err != nil {
		return fmt.Errorf("failed to create mounting. error: %w", err)
	}
	return nil
}

func (s *MountingService) CreateSeveral(ctx context.Context, mounting *mounting_api.CreateSeveralMounting) error {
	if err := s.repo.CreateSeveral(ctx, mounting); err != nil {
		return fmt.Errorf("failed to create several mountings. error: %w", err)
	}
	return nil
}

func (s *MountingService) Update(ctx context.Context, mounting *mounting_api.UpdateMounting) error {
	if err := s.repo.Update(ctx, mounting); err != nil {
		return fmt.Errorf("failed to update mounting. error: %w", err)
	}
	return nil
}

func (s *MountingService) Delete(ctx context.Context, mounting *mounting_api.DeleteMounting) error {
	if err := s.repo.Delete(ctx, mounting); err != nil {
		return fmt.Errorf("failed to delete mounting. error: %w", err)
	}
	return nil
}
