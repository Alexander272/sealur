package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_construction_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_construction_api"
)

type PutgConstructionService struct {
	repo repository.PutgConstruction
}

func NewPutgConstructionService(repo repository.PutgConstruction) *PutgConstructionService {
	return &PutgConstructionService{
		repo: repo,
	}
}

func (s *PutgConstructionService) Get(ctx context.Context, req *putg_construction_api.GetPutgConstruction) ([]*putg_construction_type_model.PutgConstruction, error) {
	constructions, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg construction. error: %w", err)
	}
	return constructions, nil
}
