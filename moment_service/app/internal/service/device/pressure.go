package device

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

func (s *DeviceService) GetPressure(ctx context.Context, req *device_api.GetPressureRequest) (pressure []*device_model.Pressure, err error) {
	data, err := s.repo.GetPressure(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get pressure. error: %w", err)
	}

	for _, item := range data {
		pressure = append(pressure, &device_model.Pressure{
			Id:    item.Id,
			Value: math.Round(item.Value*1000) / 1000,
		})
	}
	return pressure, nil
}

func (s *DeviceService) CreatePressure(ctx context.Context, pres *device_api.CreatePressureRequest) (id string, err error) {
	id, err = s.repo.CreatePressure(ctx, pres)
	if err != nil {
		return "", fmt.Errorf("failed to create pressure. error: %w", err)
	}
	return id, nil
}

func (s *DeviceService) CreateFewPressure(ctx context.Context, pres *device_api.CreateFewPressureRequest) error {
	if err := s.repo.CreateFewPressure(ctx, pres); err != nil {
		return fmt.Errorf("failed to create few pressure. error: %w", err)
	}
	return nil
}

func (s *DeviceService) UpdatePressure(ctx context.Context, pres *device_api.UpdatePressureRequest) error {
	if err := s.repo.UpdatePressure(ctx, pres); err != nil {
		return fmt.Errorf("failed to update pressure. error: %w", err)
	}
	return nil
}

func (s *DeviceService) DeletePressure(ctx context.Context, pres *device_api.DeletePressureRequest) error {
	if err := s.repo.DeletePressure(ctx, pres); err != nil {
		return fmt.Errorf("failed to delete pressure. error: %w", err)
	}
	return nil
}
