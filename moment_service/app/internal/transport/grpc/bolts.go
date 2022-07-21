package grpc

import (
	"context"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (h *FlangeHandlers) GetBolts(ctx context.Context, req *moment_proto.GetBoltsRequest) (*moment_proto.BoltsResponse, error) {
	bolts, err := h.service.GetBolts(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.BoltsResponse{Bolts: bolts}, nil
}

func (h *FlangeHandlers) CreateBolt(ctx context.Context, bolt *moment_proto.CreateBoltRequest) (*moment_proto.Response, error) {
	if err := h.service.CreateBolt(ctx, bolt); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *FlangeHandlers) UpdateBolt(ctx context.Context, bolt *moment_proto.UpdateBoltRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateBolt(ctx, bolt); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *FlangeHandlers) DeleteBolt(ctx context.Context, bolt *moment_proto.DeleteBoltRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteBolt(ctx, bolt); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
