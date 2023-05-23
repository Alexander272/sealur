package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_api"
)

type PutgFillerService struct {
	repo repository.PutgFiller
}

func NewPutgFillerService(repo repository.PutgFiller) *PutgFillerService {
	return &PutgFillerService{
		repo: repo,
	}
}

func (s *PutgFillerService) Get(ctx context.Context, req *putg_filler_api.GetPutgFiller) ([]*putg_filler_model.PutgFiller, error) {
	fillers, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg fillers. error: %w", err)
	}
	return fillers, nil
}

func (s *PutgFillerService) GetNew(ctx context.Context, req *putg_filler_api.GetPutgFiller_New) ([]*putg_filler_model.PutgFiller, error) {
	fillers, err := s.repo.GetNew(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg fillers. error: %w", err)
	}
	return fillers, nil
}

func (s *PutgFillerService) Create(ctx context.Context, filler *putg_filler_api.CreatePutgFiller) error {
	if err := s.repo.Create(ctx, filler); err != nil {
		return fmt.Errorf("failed to create putg filler. error: %w", err)
	}
	return nil
}

func (s *PutgFillerService) Update(ctx context.Context, filler *putg_filler_api.UpdatePutgFiller) error {
	if err := s.repo.Update(ctx, filler); err != nil {
		return fmt.Errorf("failed to update putg filler. error: %w", err)
	}
	return nil
}

func (s *PutgFillerService) Delete(ctx context.Context, filler *putg_filler_api.DeletePutgFiller) error {
	if err := s.repo.Delete(ctx, filler); err != nil {
		return fmt.Errorf("failed to delete putg filler. error: %w", err)
	}
	return nil
}
