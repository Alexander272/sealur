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

func (s *PutgDataService) GetByConstruction(ctx context.Context, req *putg_data_api.GetPutgData) ([]*putg_data_model.PutgData, error) {
	putgData, err := s.repo.GetByConstruction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg data. error: %w", err)
	}
	return putgData, nil
}
