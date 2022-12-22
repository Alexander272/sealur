package device

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

func (s *DeviceService) GetNumberOfMoves(ctx context.Context, req *device_api.GetNumberOfMovesRequest) (number []*device_model.NumberOfMoves, err error) {
	data, err := s.repo.GetNumberOfMoves(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get number of moves. error: %w", err)
	}

	for _, item := range data {
		number = append(number, &device_model.NumberOfMoves{
			Id:      item.Id,
			DevId:   item.DevId,
			CountId: item.CountId,
			Value:   item.Value,
		})
	}
	return number, nil
}

func (s *DeviceService) CreateNumberOfMoves(ctx context.Context, number *device_api.CreateNumberOfMovesRequest) (id string, err error) {
	id, err = s.repo.CreateNumberOfMoves(ctx, number)
	if err != nil {
		return "", fmt.Errorf("failed to create number of moves. error: %w", err)
	}
	return id, err
}

func (s *DeviceService) CreateFewNumberOfMoves(ctx context.Context, number *device_api.CreateFewNumberOfMovesRequest) error {
	if err := s.repo.CreateFewNumberOfMoves(ctx, number); err != nil {
		return fmt.Errorf("failed to create few number of moves. error :%w", err)
	}
	return nil
}

func (s *DeviceService) UpdateNumberOfMoves(ctx context.Context, number *device_api.UpdateNumberOfMovesRequest) error {
	if err := s.repo.UpdateNumberOfMoves(ctx, number); err != nil {
		return fmt.Errorf("failed to update number of moves. error: %w", err)
	}
	return nil
}

func (s *DeviceService) DeleteNumberOfMoves(ctx context.Context, number *device_api.DeleteNumberOfMovesRequest) error {
	if err := s.repo.DeleteNumberOfMoves(ctx, number); err != nil {
		return fmt.Errorf("failed to delete number of moves. error: %w", err)
	}
	return nil
}
