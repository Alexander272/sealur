package gasket

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/gasket_model"
)

func (s *GasketService) GetTypeGasket(ctx context.Context, req *gasket_api.GetGasketTypeRequest) (gasket []*gasket_model.GasketType, err error) {
	data, err := s.repo.GetTypeGasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get type gasket. error: %w", err)
	}

	for _, item := range data {
		gasket = append(gasket, &gasket_model.GasketType{
			Id:    item.Id,
			Title: item.Title,
			Label: item.Label,
		})
	}

	return gasket, nil
}

func (s *GasketService) CreateTypeGasket(ctx context.Context, gasket *gasket_api.CreateGasketTypeRequest) (id string, err error) {
	id, err = s.repo.CreateTypeGasket(ctx, gasket)
	if err != nil {
		return "", fmt.Errorf("failed to create type gasket. error: %w", err)
	}
	return id, nil
}

func (s *GasketService) UpdateTypeGasket(ctx context.Context, gasket *gasket_api.UpdateGasketTypeRequest) error {
	if err := s.repo.UpdateTypeGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to update type gasket. error: %w", err)
	}
	return nil
}

func (s *GasketService) DeleteTypeGasket(ctx context.Context, gasket *gasket_api.DeleteGasketTypeRequest) error {
	if err := s.repo.DeleteTypeGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to delete type gasket. error: %w", err)
	}
	return nil
}
