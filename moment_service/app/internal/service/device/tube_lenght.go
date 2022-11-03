package device

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

func (s *DeviceService) GetTubeLength(ctx context.Context, req *device_api.GetTubeLenghtRequest) (tube []*device_model.TubeLenght, err error) {
	data, err := s.repo.GetTubeLenght(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tube length. error: %w", err)
	}

	for _, item := range data {
		tube = append(tube, &device_model.TubeLenght{
			Id:    item.Id,
			DevId: item.DevId,
			Value: item.Value,
		})
	}
	return tube, nil
}

func (s *DeviceService) CreateTubeLength(ctx context.Context, tube *device_api.CreateTubeLenghtRequest) (id string, err error) {
	id, err = s.repo.CreateTubeLenght(ctx, tube)
	if err != nil {
		return "", fmt.Errorf("failed to create tube length. error: %w", err)
	}
	return id, err
}

func (s *DeviceService) CreateFewTubeLength(ctx context.Context, tube *device_api.CreateFewTubeLenghtRequest) error {
	if err := s.repo.CreateFewTubeLenght(ctx, tube); err != nil {
		return fmt.Errorf("failed to create few tube length")
	}
	return nil
}

func (s *DeviceService) UpdateTubeLength(ctx context.Context, tube *device_api.UpdateTubeLenghtRequest) error {
	if err := s.repo.UpdateTubeLenght(ctx, tube); err != nil {
		return fmt.Errorf("failed to update tube length. error: %w", err)
	}
	return nil
}

func (s *DeviceService) DeleteTubeLength(ctx context.Context, tube *device_api.DeleteTubeLenghtRequest) error {
	if err := s.repo.DeleteTubeLenght(ctx, tube); err != nil {
		return fmt.Errorf("failed to delete tube length. error: %w", err)
	}
	return nil
}
