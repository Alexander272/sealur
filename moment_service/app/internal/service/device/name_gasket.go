package device

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

func (s *DeviceService) GetNameGasket(ctx context.Context, req *device_api.GetNameGasketRequest) (gasket []*device_model.NameGasket, err error) {
	// data, err := s.repo.GetNameGasket(ctx, req)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get name gasket. error: %w", err)
	// }

	// for _, item := range data {
	// 	gasket = append(gasket, &device_model.NameGasket{})
	// }

	return gasket, nil
}

func (s *DeviceService) CreateNameGasket(ctx context.Context, gasket *device_api.CreateNameGasketRequest) (id string, err error) {
	id, err = s.repo.CreateNameGasket(ctx, gasket)
	if err != nil {
		return "", fmt.Errorf("failed to create name gasket. error: %w", err)
	}
	return id, nil
}

func (s *DeviceService) CreateFewNameGasket(ctx context.Context, gasket *device_api.CreateFewNameGasketRequest) error {
	if err := s.repo.CreateFewNameGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to create few name gasket. error: %w", err)
	}
	return nil
}

func (s *DeviceService) UpdateNameGasket(ctx context.Context, gasket *device_api.UpdateNameGasketRequest) error {
	if err := s.repo.UpdateNameGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to update name gasket. error: %w", err)
	}
	return nil
}

func (s *DeviceService) DeleteNameGsket(ctx context.Context, gasket *device_api.DeleteNameGasketRequest) error {
	if err := s.repo.DeleteNameGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to delete name gasker. error: %w", err)
	}
	return nil
}
