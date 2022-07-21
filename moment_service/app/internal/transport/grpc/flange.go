package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type FlangeHandlers struct {
	service service.Flange
}

func NewFlangeHandlers(service service.Flange) *FlangeHandlers {
	return &FlangeHandlers{
		service: service,
	}
}

func (h *FlangeHandlers) CreateFlangeSize(ctx context.Context, size *moment_proto.CreateFlangeSizeRequest) (*moment_proto.Response, error) {
	if err := h.service.CreateFlangeSize(ctx, size); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *FlangeHandlers) UpdateFlangeSize(ctx context.Context, size *moment_proto.UpdateFlangeSizeRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateFlangeSize(ctx, size); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *FlangeHandlers) DeleteFlangeSize(ctx context.Context, size *moment_proto.DeleteFlangeSizeRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteFlangeSize(ctx, size); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
