package materials

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/material_api"
)

func (s *MaterialsService) CreateVoltage(ctx context.Context, voltage *material_api.CreateVoltageRequest) error {
	if err := s.repo.CreateVoltage(ctx, voltage); err != nil {
		return fmt.Errorf("failed to create voltage. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) UpdateVoltage(ctx context.Context, voltage *material_api.UpdateVoltageRequest) error {
	if err := s.repo.UpdateVoltage(ctx, voltage); err != nil {
		return fmt.Errorf("failed to update voltage. error: %w", err)
	}
	return nil
}

func (s *MaterialsService) DeleteVoltage(ctx context.Context, voltage *material_api.DeleteVoltageRequest) error {
	if err := s.repo.DeleteVoltage(ctx, voltage); err != nil {
		return fmt.Errorf("failed to delete voltage. error: %w", err)
	}
	return nil
}
