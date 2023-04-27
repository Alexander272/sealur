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
