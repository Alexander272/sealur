package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (h *DeviceHandlers) GetTubeCount(ctx context.Context, req *device_api.GetTubeCountRequest) (*device_api.TubeCountResponse, error) {
	tube, err := h.service.GetTubeCount(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.TubeCountResponse{TubeCount: tube}, nil
}

func (h *DeviceHandlers) CreateTubeCount(ctx context.Context, tube *device_api.CreateTubeCountRequest) (*device_api.IdResponse, error) {
	id, err := h.service.CreateTubeCount(ctx, tube)
	if err != nil {
		return nil, err
	}
	return &device_api.IdResponse{Id: id}, nil
}

func (h *DeviceHandlers) CreateFewTubeCount(ctx context.Context, tube *device_api.CreateFewTubeCountRequest) (*device_api.Response, error) {
	if err := h.service.CreateFewTubeCount(ctx, tube); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) UpdateTubeCount(ctx context.Context, tube *device_api.UpdateTubeCountRequest) (*device_api.Response, error) {
	if err := h.service.UpdateTubeCount(ctx, tube); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) DeleteTubeCount(ctx context.Context, tube *device_api.DeleteTubeCountRequest) (*device_api.Response, error) {
	if err := h.service.DeleteTubeCount(ctx, tube); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
