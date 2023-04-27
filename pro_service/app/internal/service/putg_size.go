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
