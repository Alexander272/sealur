package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_configuration_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_conf_api"
)

type PutgConfigurationService struct {
	repo repository.PutgConfiguration
}

func NewPutgConfigurationService(repo repository.PutgConfiguration) *PutgConfigurationService {
	return &PutgConfigurationService{
		repo: repo,
	}
}

func (s *PutgConfigurationService) Get(ctx context.Context, req *putg_conf_api.GetPutgConfiguration) ([]*putg_configuration_model.PutgConfiguration, error) {
	configurations, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg configuration. error: %w", err)
	}
	return configurations, nil
}
