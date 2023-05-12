package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_data_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_data_api"
)

type PutgDataService struct {
	repo repository.PutgData
}

func NewPutgDataService(repo repository.PutgData) *PutgDataService {
	return &PutgDataService{
		repo: repo,
	}
}

func (s *PutgDataService) Get(ctx context.Context, req *putg_data_api.GetPutgData) (*putg_data_model.PutgData, error) {
	putgData, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg data. error: %w", err)
	}
	return putgData, nil
}

func (s *PutgDataService) Create(ctx context.Context, data *putg_data_api.CreatePutgData) error {
	if err := s.repo.Create(ctx, data); err != nil {
		return fmt.Errorf("failed to create putg data. error: %w", err)
	}
	return nil
}

func (s *PutgDataService) Update(ctx context.Context, data *putg_data_api.UpdatePutgData) error {
	if err := s.repo.Update(ctx, data); err != nil {
		return fmt.Errorf("failed to update putg data. error: %w", err)
	}
	return nil
}

func (s *PutgDataService) Delete(ctx context.Context, data *putg_data_api.DeletePutgData) error {
	if err := s.repo.Delete(ctx, data); err != nil {
		return fmt.Errorf("failed to delete putg data. error: %w", err)
	}
	return nil
}
