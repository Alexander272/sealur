package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (h *DeviceHandlers) GetNumberOfMoves(ctx context.Context, req *device_api.GetNumberOfMovesRequest) (*device_api.NumberOfMovesResponse, error) {
	number, err := h.service.GetNumberOfMoves(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_api.NumberOfMovesResponse{Number: number}, nil
}

func (h *DeviceHandlers) CreateNumberOfMoves(ctx context.Context, number *device_api.CreateNumberOfMovesRequest) (*device_api.IdResponse, error) {
	id, err := h.service.CreateNumberOfMoves(ctx, number)
	if err != nil {
		return nil, err
	}
	return &device_api.IdResponse{Id: id}, nil
}

func (h *DeviceHandlers) CreateFewNumberOfMoves(ctx context.Context, number *device_api.CreateFewNumberOfMovesRequest) (*device_api.Response, error) {
	if err := h.service.CreateFewNumberOfMoves(ctx, number); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) UpdateNumberOfMoves(ctx context.Context, number *device_api.UpdateNumberOfMovesRequest) (*device_api.Response, error) {
	if err := h.service.UpdateNumberOfMoves(ctx, number); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}

func (h *DeviceHandlers) DeleteNumberOfMoves(ctx context.Context, number *device_api.DeleteNumberOfMovesRequest) (*device_api.Response, error) {
	if err := h.service.DeleteNumberOfMoves(ctx, number); err != nil {
		return nil, err
	}
	return &device_api.Response{}, nil
}
