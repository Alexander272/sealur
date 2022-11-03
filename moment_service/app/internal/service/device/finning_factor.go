package device

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

func (s *DeviceService) GetFinningFactor(ctx context.Context, req *device_api.GetFinningFactorRequest) (factor []*device_model.FinningFactor, err error) {
	data, err := s.repo.GetFinningFactor(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get finning factor. error: %w", err)
	}

	for _, item := range data {
		factor = append(factor, &device_model.FinningFactor{
			Id:    item.Id,
			DevId: item.DevId,
			Value: item.Value,
		})
	}
	return factor, nil
}

func (s *DeviceService) CreateFinningFactor(ctx context.Context, factor *device_api.CreateFinningFactorRequest) (id string, err error) {
	id, err = s.repo.CreateFinningFactor(ctx, factor)
	if err != nil {
		return "", fmt.Errorf("failed to create finning factor. error: %w", err)
	}
	return id, nil
}

func (s *DeviceService) CreateFewFinningFactor(ctx context.Context, factor *device_api.CreateFewFinningFactorRequest) error {
	if err := s.repo.CreateFewFinningFactor(ctx, factor); err != nil {
		return fmt.Errorf("failed to create few finning factor. error: %w", err)
	}
	return nil
}

func (s *DeviceService) UpdateFinningFactor(ctx context.Context, factor *device_api.UpdateFinningFactorRequest) error {
	if err := s.repo.UpdateFinningFactor(ctx, factor); err != nil {
		return fmt.Errorf("failed to update finning factor. error: %w", err)
	}
	return nil
}

func (s *DeviceService) DeleteFinningFactor(ctx context.Context, factor *device_api.DeleteFinningFactorRequest) error {
	if err := s.repo.DeleteFinningFactor(ctx, factor); err != nil {
		return fmt.Errorf("failed to delete finning factor. error: %w", err)
	}
	return nil
}
