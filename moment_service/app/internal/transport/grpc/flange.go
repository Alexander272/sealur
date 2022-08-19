package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type FlangeHandlers struct {
	service service.Flange
	moment_api.UnimplementedFlangeServiceServer
}

func NewFlangeHandlers(service service.Flange) *FlangeHandlers {
	return &FlangeHandlers{
		service: service,
	}
}

func (h *FlangeHandlers) GetFlangeSize(ctx context.Context, size *moment_api.GetFullFlangeSizeRequest) (*moment_api.FullFlangeSizeResponse, error) {
	sizes, err := h.service.GetFullFlangeSize(ctx, size)
	if err != nil {
		return nil, err
	}
	return sizes, nil
}

func (h *FlangeHandlers) CreateFlangeSize(ctx context.Context, size *moment_api.CreateFlangeSizeRequest) (*moment_api.Response, error) {
	if err := h.service.CreateFlangeSize(ctx, size); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *FlangeHandlers) UpdateFlangeSize(ctx context.Context, size *moment_api.UpdateFlangeSizeRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateFlangeSize(ctx, size); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *FlangeHandlers) DeleteFlangeSize(ctx context.Context, size *moment_api.DeleteFlangeSizeRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteFlangeSize(ctx, size); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}
