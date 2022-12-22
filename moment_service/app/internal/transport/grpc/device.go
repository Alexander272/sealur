package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

type DeviceHandlers struct {
	service service.Device
	device_api.UnimplementedDeviceServiceServer
}

func NewDeviceHandlers(service service.Device) *DeviceHandlers {
	return &DeviceHandlers{
		service: service,
	}
}

func (h *DeviceHandlers) GetDevice(ctx context.Context, req *device_api.GetDeviceRequest) (*device_api.DeviceResponse, error) {
	devices, err := h.service.GetDevices(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.DeviceResponse{Devices: devices}, nil
}

func (h *DeviceHandlers) CreateDevice(ctx context.Context, device *device_api.CreateDeviceRequest) (*device_api.IdResponse, error) {
	id, err := h.service.CreateDevice(ctx, device)
	if err != nil {
		return nil, err
	}
	return &device_api.IdResponse{Id: id}, nil
}

func (h *DeviceHandlers) CreateFewDevice(ctx context.Context, devices *device_api.CreateFewDeviceRequest) (*device_api.Response, error) {
	if err := h.service.CreateFewDevices(ctx, devices); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) UpdateDevice(ctx context.Context, device *device_api.UpdateDeviceRequest) (*device_api.Response, error) {
	if err := h.service.UpdateDevice(ctx, device); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) DeleteDevice(ctx context.Context, device *device_api.DeleteDeviceRequest) (*device_api.Response, error) {
	if err := h.service.DeleteDevice(ctx, device); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
