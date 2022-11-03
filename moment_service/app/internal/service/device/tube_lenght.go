package device

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

func (s *DeviceService) GetTubeLength(ctx context.Context, req *device_api.GetTubeLengthRequest) (tube []*device_model.TubeLength, err error) {
	data, err := s.repo.GetTubeLength(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tube length. error: %w", err)
	}

	for _, item := range data {
		tube = append(tube, &device_model.TubeLength{
			Id:    item.Id,
			DevId: item.DevId,
			Value: item.Value,
		})
	}
	return tube, nil
}

func (s *DeviceService) CreateTubeLength(ctx context.Context, tube *device_api.CreateTubeLengthRequest) (id string, err error) {
	id, err = s.repo.CreateTubeLength(ctx, tube)
	if err != nil {
		return "", fmt.Errorf("failed to create tube length. error: %w", err)
	}
	return id, err
}

func (s *DeviceService) CreateFewTubeLength(ctx context.Context, tube *device_api.CreateFewTubeLengthRequest) error {
	if err := s.repo.CreateFewTubeLength(ctx, tube); err != nil {
		return fmt.Errorf("failed to create few tube length")
	}
	return nil
}

func (s *DeviceService) UpdateTubeLength(ctx context.Context, tube *device_api.UpdateTubeLengthRequest) error {
	if err := s.repo.UpdateTubeLength(ctx, tube); err != nil {
		return fmt.Errorf("failed to update tube length. error: %w", err)
	}
	return nil
}

func (s *DeviceService) DeleteTubeLength(ctx context.Context, tube *device_api.DeleteTubeLengthRequest) error {
	if err := s.repo.DeleteTubeLength(ctx, tube); err != nil {
		return fmt.Errorf("failed to delete tube length. error: %w", err)
	}
	return nil
}
