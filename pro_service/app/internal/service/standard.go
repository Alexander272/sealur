package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/standard_api"
)

type StandardService struct {
	repo repository.Standard
}

func NewStandardService(repo repository.Standard) *StandardService {
	return &StandardService{repo: repo}
}

func (s *StandardService) GetAll(ctx context.Context, standard *standard_api.GetAllStandards) ([]*standard_model.Standard, error) {
	standards, err := s.repo.GetAll(ctx, standard)
	if err != nil {
		return nil, fmt.Errorf("failed to get standard. error: %w", err)
	}
	return standards, err
}

func (s *StandardService) Create(ctx context.Context, standard *standard_api.CreateStandard) error {
	if err := s.repo.Create(ctx, standard); err != nil {
		return fmt.Errorf("failed to create standard. error: %w", err)
	}
	return nil
}

func (s *StandardService) CreateSeveral(ctx context.Context, standards *standard_api.CreateSeveralStandard) error {
	if err := s.repo.CreateSeveral(ctx, standards); err != nil {
		return fmt.Errorf("failed to create several standard standards. error: %w", err)
	}
	return nil
}

func (s *StandardService) Update(ctx context.Context, standard *standard_api.UpdateStandard) error {
	if err := s.repo.Update(ctx, standard); err != nil {
		return fmt.Errorf("failed to update standard standard. error: %w", err)
	}
	return nil
}

func (s *StandardService) Delete(ctx context.Context, standard *standard_api.DeleteStandard) error {
	if err := s.repo.Delete(ctx, standard); err != nil {
		return fmt.Errorf("failed to delete standard standard. error: %w", err)
	}
	return nil
}
