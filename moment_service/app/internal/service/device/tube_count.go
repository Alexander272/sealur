package device

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

func (s *DeviceService) GetTubeCount(ctx context.Context, req *device_api.GetTubeCountRequest) (tube []*device_model.TubeCount, err error) {
	data, err := s.repo.GetTubeCount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tube count. error: %w", err)
	}

	for _, item := range data {
		tube = append(tube, &device_model.TubeCount{
			Id:    item.Id,
			Value: item.Value,
		})
	}
	return tube, nil
}

func (s *DeviceService) CreateTubeCount(ctx context.Context, tube *device_api.CreateTubeCountRequest) (id string, err error) {
	id, err = s.repo.CreateTubeCount(ctx, tube)
	if err != nil {
		return "", fmt.Errorf("failed to create tube count. error: %w", err)
	}
	return id, nil
}

func (s *DeviceService) CreateFewTubeCount(ctx context.Context, tube *device_api.CreateFewTubeCountRequest) error {
	if err := s.repo.CreateFewTubeCount(ctx, tube); err != nil {
		return fmt.Errorf("failed to create few tube count. error: %w", err)
	}
	return nil
}

func (s *DeviceService) UpdateTubeCount(ctx context.Context, tube *device_api.UpdateTubeCountRequest) error {
	if err := s.repo.UpdateTubeCount(ctx, tube); err != nil {
		return fmt.Errorf("failed to update tube count. error: %w", err)
	}
	return nil
}

func (s *DeviceService) DeleteTubeCount(ctx context.Context, tube *device_api.DeleteTubeCountRequest) error {
	if err := s.repo.DeleteTubeCount(ctx, tube); err != nil {
		return fmt.Errorf("failed to delete tube count. error: %w", err)
	}
	return nil
}
