package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_standard_api"
)

type PutgStandardService struct {
	repo repository.PutgStandard
}

func NewPutgStandardService(repo repository.PutgStandard) *PutgStandardService {
	return &PutgStandardService{
		repo: repo,
	}
}

func (s *PutgStandardService) Get(ctx context.Context, req *putg_standard_api.GetPutgStandard) ([]*putg_standard_model.PutgStandard, error) {
	standards, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg standards. error: %w", err)
	}
	return standards, nil
}

func (s *PutgStandardService) Create(ctx context.Context, standard *putg_standard_api.CreatePutgStandard) error {
	if err := s.repo.Create(ctx, standard); err != nil {
		return fmt.Errorf("failed to create putg standard. error: %w", err)
	}
	return nil
}

func (s *PutgStandardService) Update(ctx context.Context, standard *putg_standard_api.UpdatePutgStandard) error {
	if err := s.repo.Update(ctx, standard); err != nil {
		return fmt.Errorf("failed to update putg standard. error: %w", err)
	}
	return nil
}

func (s *PutgStandardService) Delete(ctx context.Context, standard *putg_standard_api.DeletePutgStandard) error {
	if err := s.repo.Delete(ctx, standard); err != nil {
		return fmt.Errorf("failed to delete putg standard. error: %w", err)
	}
	return nil
}
