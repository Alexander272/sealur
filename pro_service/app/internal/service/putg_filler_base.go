package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_base_api"
)

type PutgBaseFillerService struct {
	repo repository.PutgBaseFiller
}

func NewPutgBaseFillerService(repo repository.PutgBaseFiller) *PutgBaseFillerService {
	return &PutgBaseFillerService{
		repo: repo,
	}
}

func (s *PutgBaseFillerService) Get(ctx context.Context, req *putg_filler_base_api.GetPutgBaseFiller) ([]*putg_filler_model.PutgFiller, error) {
	fillers, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg base fillers. error: %w", err)
	}

	return fillers, nil
}

func (s *PutgBaseFillerService) Create(ctx context.Context, filler *putg_filler_base_api.CreatePutgBaseFiller) error {
	if err := s.repo.Create(ctx, filler); err != nil {
		return fmt.Errorf("failed to create putg base filler. error: %w", err)
	}
	return nil
}

func (s *PutgBaseFillerService) Update(ctx context.Context, filler *putg_filler_base_api.UpdatePutgBaseFiller) error {
	if err := s.repo.Update(ctx, filler); err != nil {
		return fmt.Errorf("failed to update putg base filler. error: %w", err)
	}
	return nil
}

func (s *PutgBaseFillerService) Delete(ctx context.Context, filler *putg_filler_base_api.DeletePutgBaseFiller) error {
	if err := s.repo.Delete(ctx, filler); err != nil {
		return fmt.Errorf("failed to delete putg base filler. error: %w", err)
	}
	return nil
}
