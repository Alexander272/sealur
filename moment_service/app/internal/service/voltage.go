package service

import (
	"context"
	"fmt"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (s *MaterialsService) CreateVoltage(ctx context.Context, voltage *moment_proto.CreateVoltageRequest) error {
	if err := s.repo.CreateVoltage(ctx, voltage); err != nil {
		return fmt.Errorf("failed to create voltage. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) UpdateVoltage(ctx context.Context, voltage *moment_proto.UpdateVoltageRequest) error {
	if err := s.repo.UpdateVoltage(ctx, voltage); err != nil {
		return fmt.Errorf("failed to update voltage. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteVoltage(ctx context.Context, voltage *moment_proto.DeleteVoltageRequest) error {
	if err := s.repo.DeleteVoltage(ctx, voltage); err != nil {
		return fmt.Errorf("failed to delete voltage. error: %w", err)
	}
	return nil
}
