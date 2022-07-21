package grpc

import (
	"context"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (h *MaterialsHandlers) CreateAlpha(ctx context.Context, alpha *moment_proto.CreateAlphaRequest) (*moment_proto.Response, error) {
	err := h.service.CreateAlpha(ctx, alpha)
	if err != nil {
		return nil, err
	}

	return &moment_proto.Response{}, nil
}

func (h *MaterialsHandlers) UpdateAlpha(ctx context.Context, alpha *moment_proto.UpdateAlphaRequest) (*moment_proto.Response, error) {
	if err := h.service.UpateAlpha(ctx, alpha); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *MaterialsHandlers) DeleteAlpha(ctx context.Context, alpha *moment_proto.DeleteAlphaRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteAlpha(ctx, alpha); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
