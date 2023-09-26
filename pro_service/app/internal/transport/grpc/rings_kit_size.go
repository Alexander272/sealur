package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_size_api"
)

type RingsKitSizeHandlers struct {
	service service.RingsKitSize
	rings_kit_size_api.UnimplementedRingSizesKitServiceServer
}

func NewRingsKitSizeHandlers(service service.RingsKitSize) *RingsKitSizeHandlers {
	return &RingsKitSizeHandlers{
		service: service,
	}
}

func (h *RingsKitSizeHandlers) GetAll(ctx context.Context, req *rings_kit_size_api.GetRingsKitSize) (*rings_kit_size_api.RingsKitSize, error) {
	sizes, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &rings_kit_size_api.RingsKitSize{Sizes: sizes}, nil
}

func (h *RingsKitSizeHandlers) Create(ctx context.Context, size *rings_kit_size_api.CreateRingsKitSize) (*response_model.Response, error) {
	if err := h.service.Create(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingsKitSizeHandlers) Update(ctx context.Context, size *rings_kit_size_api.UpdateRingsKitSize) (*response_model.Response, error) {
	if err := h.service.Update(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingsKitSizeHandlers) Delete(ctx context.Context, size *rings_kit_size_api.DeleteRingsKitSize) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
