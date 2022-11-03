package device

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/moment_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

type DeviceService struct {
	repo repository.Device
}

func NewDeviceService(repo repository.Device) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) GetDevices(ctx context.Context, req *device_api.GetDeviceRequest) (devices []*device_model.Device, err error) {
	data, err := s.repo.GetDevices(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices. error: %w", err)
	}

	for _, item := range data {
		devices = append(devices, &device_model.Device{
			Id:    item.Id,
			Title: item.Title,
		})
	}
	return devices, nil
}

func (s *DeviceService) CreateDevice(ctx context.Context, device *device_api.CreateDeviceRequest) (id string, err error) {
	id, err = s.repo.CreateDevice(ctx, device)
	if err != nil {
		return "", fmt.Errorf("failed to create device. error: %w", err)
	}

	return id, err
}

func (s *DeviceService) CreateFewDevices(ctx context.Context, device *device_api.CreateFewDeviceRequest) error {
	if err := s.repo.CreateFewDevices(ctx, device); err != nil {
		return fmt.Errorf("failed to create few devices. error: %w", err)
	}
	return nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, device *device_api.UpdateDeviceRequest) error {
	if err := s.repo.UpdateDevice(ctx, device); err != nil {
		return fmt.Errorf("failed to update device. error: %w", err)
	}
	return nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, device *device_api.DeleteDeviceRequest) error {
	if err := s.repo.DeleteDevice(ctx, device); err != nil {
		return fmt.Errorf("failed to delete device. error: %w", err)
	}
	return nil
}
