package gasket

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/gasket_model"
)

func (s *GasketService) GetEnv(ctx context.Context, req *gasket_api.GetEnvRequest) (env []*gasket_model.Env, err error) {
	data, err := s.repo.GetEnv(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get env. error: %w", err)
	}

	for _, item := range data {
		env = append(env, &gasket_model.Env{
			Id:    item.Id,
			Title: item.Title,
		})
	}

	return env, nil
}

func (s *GasketService) CreateEnv(ctx context.Context, env *gasket_api.CreateEnvRequest) (id string, err error) {
	id, err = s.repo.CreateEnv(ctx, env)
	if err != nil {
		return "", fmt.Errorf("failed to create env. error: %w", err)
	}
	return id, nil
}

func (s *GasketService) UpdateEnv(ctx context.Context, env *gasket_api.UpdateEnvRequest) error {
	if err := s.repo.UpdateEnv(ctx, env); err != nil {
		return fmt.Errorf("failed to update env. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteEnv(ctx context.Context, env *gasket_api.DeleteEnvRequest) error {
	if err := s.repo.DeleteEnv(ctx, env); err != nil {
		return fmt.Errorf("failed to delete env. error: %w", err)
	}
	return nil
}

//---
func (s *GasketService) CreateManyEnvData(ctx context.Context, data *gasket_api.CreateManyEnvDataRequest) error {
	if err := s.repo.CreateManyEnvData(ctx, data); err != nil {
		return fmt.Errorf("failed to create many env data. error: %w", err)
	}
	return nil
}

func (s *GasketService) CreateEnvData(ctx context.Context, data *gasket_api.CreateEnvDataRequest) error {
	if err := s.repo.CreateEnvData(ctx, data); err != nil {
		return fmt.Errorf("failed to create env data. error: %w", err)
	}
	return nil
}

func (s *GasketService) UpdateEnvData(ctx context.Context, data *gasket_api.UpdateEnvDataRequest) error {
	if err := s.repo.UpdateEnvData(ctx, data); err != nil {
		return fmt.Errorf("failed to update env data. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteEnvData(ctx context.Context, data *gasket_api.DeleteEnvDataRequest) error {
	if err := s.repo.DeleteEnvData(ctx, data); err != nil {
		return fmt.Errorf("failed to delete env data. error: %w", err)
	}
	return nil
}
