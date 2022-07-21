package service

import (
	"context"
	"fmt"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

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
