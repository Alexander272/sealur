package grpc

import (
	"context"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (h *MaterialsHandlers) CreateVoltage(ctx context.Context, voltage *moment_proto.CreateVoltageRequest) (*moment_proto.Response, error) {
	err := h.service.CreateVoltage(ctx, voltage)
	if err != nil {
		return nil, err
	}

	return &moment_proto.Response{}, nil
}

func (h *MaterialsHandlers) UpdateVoltage(ctx context.Context, voltage *moment_proto.UpdateVoltageRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateVoltage(ctx, voltage); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *MaterialsHandlers) DeleteVoltage(ctx context.Context, voltage *moment_proto.DeleteVoltageRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteVoltage(ctx, voltage); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
