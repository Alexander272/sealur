package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (h *DeviceHandlers) GetTubeLenght(ctx context.Context, req *device_api.GetTubeLenghtRequest) (*device_api.TubeLenghtResponse, error) {
	tube, err := h.service.GetTubeLength(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.TubeLenghtResponse{TubeLenght: tube}, nil
}

func (h *DeviceHandlers) CreateTubeLenght(ctx context.Context, tube *device_api.CreateTubeLenghtRequest) (*device_api.IdResponse, error) {
	id, err := h.service.CreateTubeLength(ctx, tube)
	if err != nil {
		return nil, err
	}
	return &device_api.IdResponse{Id: id}, nil
}

func (h *DeviceHandlers) CreateFewTubeLenght(ctx context.Context, tube *device_api.CreateFewTubeLenghtRequest) (*device_api.Response, error) {
	if err := h.service.CreateFewTubeLength(ctx, tube); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) UpdateTubeLenght(ctx context.Context, tube *device_api.UpdateTubeLenghtRequest) (*device_api.Response, error) {
	if err := h.service.UpdateTubeLength(ctx, tube); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) DeleteTubeLenght(ctx context.Context, tube *device_api.DeleteTubeLenghtRequest) (*device_api.Response, error) {
	if err := h.service.DeleteTubeLength(ctx, tube); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
