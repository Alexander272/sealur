package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type GasketService struct {
	repo repository.Gasket
}

func NewGasketService(repo repository.Gasket) *GasketService {
	return &GasketService{
		repo: repo,
	}
}

func (s *GasketService) GetFullData(ctx context.Context, gasket models.GetGasket) (models.FullDataGasket, error) {
	g, err := s.repo.GetFullData(ctx, gasket)
	if err != nil {
		return models.FullDataGasket{}, fmt.Errorf("failed to get gasket. error: %w", err)
	}
	return g, nil
}

func (s *GasketService) GetGasket(ctx context.Context, req *moment_proto.GetGasketRequest) (gasket []*moment_proto.Gasket, err error) {
	data, err := s.repo.GetGasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get gasket. error: %w", err)
	}

	for _, item := range data {
		g := moment_proto.Gasket(item)
		gasket = append(gasket, &g)
	}

	return gasket, nil
}

func (s *GasketService) CreateGasket(ctx context.Context, gasket *moment_proto.CreateGasketRequest) (id string, err error) {
	id, err = s.repo.CreateGasket(ctx, gasket)
	if err != nil {
		return "", fmt.Errorf("failed to create gasket. error: %w", err)
	}
	return id, nil
}

func (s *GasketService) UpdateGasket(ctx context.Context, gasket *moment_proto.UpdateGasketRequest) error {
	if err := s.repo.UpdateGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to update gasket. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteGasket(ctx context.Context, gasket *moment_proto.DeleteGasketRequest) error {
	if err := s.repo.DeleteGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to delete gasket. error: %w", err)
	}
	return nil
}

//---

func (s *GasketService) GetTypeGasket(ctx context.Context, req *moment_proto.GetGasketTypeRequest) (gasket []*moment_proto.GasketType, err error) {
	data, err := s.repo.GetTypeGasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get type gasket. error: %w", err)
	}

	for _, item := range data {
		g := moment_proto.GasketType(item)
		gasket = append(gasket, &g)
	}

	return gasket, nil
}

func (s *GasketService) CreateTypeGasket(ctx context.Context, gasket *moment_proto.CreateGasketTypeRequest) (id string, err error) {
	id, err = s.repo.CreateTypeGasket(ctx, gasket)
	if err != nil {
		return "", fmt.Errorf("failed to create type gasket. error: %w", err)
	}
	return id, nil
}

func (s *GasketService) UpdateTypeGasket(ctx context.Context, gasket *moment_proto.UpdateGasketTypeRequest) error {
	if err := s.repo.UpdateTypeGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to update type gasket. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteTypeGasket(ctx context.Context, gasket *moment_proto.DeleteGasketTypeRequest) error {
	if err := s.repo.DeleteTypeGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to delete type gasket. error: %w", err)
	}
	return nil
}

//---

func (s *GasketService) GetEnv(ctx context.Context, req *moment_proto.GetEnvRequest) (env []*moment_proto.Env, err error) {
	data, err := s.repo.GetEnv(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get env. error: %w", err)
	}

	for _, item := range data {
		e := moment_proto.Env(item)
		env = append(env, &e)
	}

	return env, nil
}

func (s *GasketService) CreateEnv(ctx context.Context, env *moment_proto.CreateEnvRequest) (id string, err error) {
	id, err = s.repo.CreateEnv(ctx, env)
	if err != nil {
		return "", fmt.Errorf("failed to create env. error: %w", err)
	}
	return id, nil
}

func (s *GasketService) UpdateEnv(ctx context.Context, env *moment_proto.UpdateEnvRequest) error {
	if err := s.repo.UpdateEnv(ctx, env); err != nil {
		return fmt.Errorf("failed to update env. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteEnv(ctx context.Context, env *moment_proto.DeleteEnvRequest) error {
	if err := s.repo.DeleteEnv(ctx, env); err != nil {
		return fmt.Errorf("failed to delete env. error: %w", err)
	}
	return nil
}

//---

func (s *GasketService) CreateEnvData(ctx context.Context, data *moment_proto.CreateEnvDataRequest) error {
	if err := s.repo.CreateEnvData(ctx, data); err != nil {
		return fmt.Errorf("failed to create env data. error: %w", err)
	}
	return nil
}

func (s *GasketService) UpdateEnvData(ctx context.Context, data *moment_proto.UpdateEnvDataRequest) error {
	if err := s.repo.UpdateEnvData(ctx, data); err != nil {
		return fmt.Errorf("failed to update env data. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteEnvData(ctx context.Context, data *moment_proto.DeleteEnvDataRequest) error {
	if err := s.repo.DeleteEnvData(ctx, data); err != nil {
		return fmt.Errorf("failed to delete env data. error: %w", err)
	}
	return nil
}

//---

func (s *GasketService) CreateGasketData(ctx context.Context, data *moment_proto.CreateGasketDataRequest) error {
	if err := s.repo.CreateGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to create gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) UpdateGasketData(ctx context.Context, data *moment_proto.UpdateGasketDataRequest) error {
	if err := s.repo.UpdateGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to update gasket data. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteGasketData(ctx context.Context, data *moment_proto.DeleteGasketDataRequest) error {
	if err := s.repo.DeleteGasketData(ctx, data); err != nil {
		return fmt.Errorf("failed to delete gasket data. error: %w", err)
	}
	return nil
}
