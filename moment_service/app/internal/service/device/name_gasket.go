package device

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

func (s *DeviceService) GetNameGasket(ctx context.Context, req *device_api.GetNameGasketRequest) (gasket []*device_model.NameGasket, err error) {
	data, err := s.repo.GetNameGasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get name gasket. error: %w", err)
	}

	for _, item := range data {
		gasket = append(gasket, &device_model.NameGasket{
			Id:     item.Id,
			FinId:  req.FinId,
			NumId:  item.NumId,
			PresId: item.PresId,
			Title:  item.Title,
		})
	}
	return gasket, nil
}

func (s *DeviceService) GetFullNameGasket(ctx context.Context, req *device_api.GetFullNameGasketRequest) (gasket []*device_model.FullNameGasket, err error) {
	data, err := s.repo.GetFullNameGasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get name gasket. error: %w", err)
	}

	for _, item := range data {
		gasket = append(gasket, &device_model.FullNameGasket{
			Id:        item.Id,
			FinId:     req.FinId,
			NumId:     item.NumId,
			PresId:    item.PresId,
			Title:     item.Title,
			SizeLong:  math.Round(item.SizeLong*1000) / 1000,
			SizeTrans: math.Round(item.SizeTrans*1000) / 1000,
			Width:     math.Round(item.Width*1000) / 1000,
			Thick1:    math.Round(item.Thick1*1000) / 1000,
			Thick2:    math.Round(item.Thick2*1000) / 1000,
			Thick3:    math.Round(item.Thick3*1000) / 1000,
			Thick4:    math.Round(item.Thick4*1000) / 1000,
		})
	}
	return gasket, nil
}

func (s *DeviceService) GetNameGasketSize(ctx context.Context, req *device_api.GetNameGasketSizeRequest) (gasket *device_model.NameGasketSize, err error) {
	data, err := s.repo.GetNameGasketSize(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get name gasket. error: %w", err)
	}

	gasket = &device_model.NameGasketSize{
		Id:        data.Id,
		SizeLong:  math.Round(data.SizeLong*1000) / 1000,
		SizeTrans: math.Round(data.SizeTrans*1000) / 1000,
		Width:     math.Round(data.Width*1000) / 1000,
		Thick1:    math.Round(data.Thick1*1000) / 1000,
		Thick2:    math.Round(data.Thick2*1000) / 1000,
		Thick3:    math.Round(data.Thick3*1000) / 1000,
		Thick4:    math.Round(data.Thick4*1000) / 1000,
	}
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

func (s *DeviceService) DeleteNameGasket(ctx context.Context, gasket *device_api.DeleteNameGasketRequest) error {
	if err := s.repo.DeleteNameGasket(ctx, gasket); err != nil {
		return fmt.Errorf("failed to delete name gasket. error: %w", err)
	}
	return nil
}
