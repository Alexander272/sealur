package service

import (
	"context"
	"fmt"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

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
