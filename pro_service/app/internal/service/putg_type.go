package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_type_api"
)

type PutgTypeService struct {
	repo repository.PutgType
}

func NewPutgTypeService(repo repository.PutgType) *PutgTypeService {
	return &PutgTypeService{
		repo: repo,
	}
}

func (s *PutgTypeService) Get(ctx context.Context, req *putg_type_api.GetPutgType) ([]*putg_type_model.PutgType, error) {
	types, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg types. error: %w", err)
	}
	return types, nil
}
