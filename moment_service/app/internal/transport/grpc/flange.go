package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
)

type FlangeHandlers struct {
	service service.Flange
	flange_api.UnimplementedFlangeServiceServer
}

func NewFlangeHandlers(service service.Flange) *FlangeHandlers {
	return &FlangeHandlers{
		service: service,
	}
}

func (h *FlangeHandlers) GetFlangeSize(ctx context.Context, size *flange_api.GetFullFlangeSizeRequest) (*flange_api.FullFlangeSizeResponse, error) {
	sizes, err := h.service.GetFullFlangeSize(ctx, size)
	if err != nil {
		return nil, err
	}
	return sizes, nil
}

func (h *FlangeHandlers) CreateFlangeSize(ctx context.Context, size *flange_api.CreateFlangeSizeRequest) (*flange_api.Response, error) {
	if err := h.service.CreateFlangeSize(ctx, size); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}

func (h *FlangeHandlers) CreateFlangeSizes(ctx context.Context, size *flange_api.CreateFlangeSizesRequest) (*flange_api.Response, error) {
	if err := h.service.CreateFlangeSizes(ctx, size); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}

func (h *FlangeHandlers) UpdateFlangeSize(ctx context.Context, size *flange_api.UpdateFlangeSizeRequest) (*flange_api.Response, error) {
	if err := h.service.UpdateFlangeSize(ctx, size); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}

func (h *FlangeHandlers) DeleteFlangeSize(ctx context.Context, size *flange_api.DeleteFlangeSizeRequest) (*flange_api.Response, error) {
	if err := h.service.DeleteFlangeSize(ctx, size); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}
