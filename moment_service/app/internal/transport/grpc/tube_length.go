package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (h *DeviceHandlers) GetTubeLength(ctx context.Context, req *device_api.GetTubeLengthRequest) (*device_api.TubeLengthResponse, error) {
	tube, err := h.service.GetTubeLength(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.TubeLengthResponse{TubeLength: tube}, nil
}

func (h *DeviceHandlers) CreateTubeLength(ctx context.Context, tube *device_api.CreateTubeLengthRequest) (*device_api.IdResponse, error) {
	id, err := h.service.CreateTubeLength(ctx, tube)
	if err != nil {
		return nil, err
	}
	return &device_api.IdResponse{Id: id}, nil
}

func (h *DeviceHandlers) CreateFewTubeLength(ctx context.Context, tube *device_api.CreateFewTubeLengthRequest) (*device_api.Response, error) {
	if err := h.service.CreateFewTubeLength(ctx, tube); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) UpdateTubeLength(ctx context.Context, tube *device_api.UpdateTubeLengthRequest) (*device_api.Response, error) {
	if err := h.service.UpdateTubeLength(ctx, tube); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) DeleteTubeLength(ctx context.Context, tube *device_api.DeleteTubeLengthRequest) (*device_api.Response, error) {
	if err := h.service.DeleteTubeLength(ctx, tube); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
