package grpc

import (
	"context"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (h *GasketHandlers) GetGasketType(ctx context.Context, req *moment_proto.GetGasketTypeRequest) (*moment_proto.GasketTypeResponse, error) {
	typeGasket, err := h.service.GetTypeGasket(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.GasketTypeResponse{GasketType: typeGasket}, nil
}

func (h *GasketHandlers) CreateGasketType(ctx context.Context, gasket *moment_proto.CreateGasketTypeRequest) (*moment_proto.IdResponse, error) {
	id, err := h.service.CreateTypeGasket(ctx, gasket)
	if err != nil {
		return nil, err
	}

	return &moment_proto.IdResponse{Id: id}, nil
}

func (h *GasketHandlers) UpdateGasketType(ctx context.Context, gasket *moment_proto.UpdateGasketTypeRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateTypeGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *GasketHandlers) DeleteGasketType(ctx context.Context, gasket *moment_proto.DeleteGasketTypeRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteTypeGasket(ctx, gasket); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
