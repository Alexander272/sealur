package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_size_api"
)

type RingSizeHandlers struct {
	service service.RingSize
	*ring_size_api.UnimplementedRingSizeServiceServer
}

func NewRingSizeHandlers(service service.RingSize) *RingSizeHandlers {
	return &RingSizeHandlers{
		service: service,
	}
}

func (h *RingSizeHandlers) GetAll(ctx context.Context, req *ring_size_api.GetRingSize) (*ring_size_api.RingSize, error) {
	sizes, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ring_size_api.RingSize{Size: sizes}, nil
}

func (h *RingSizeHandlers) Create(ctx context.Context, size *ring_size_api.CreateRingSize) (*response_model.Response, error) {
	if err := h.service.Create(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingSizeHandlers) Update(ctx context.Context, size *ring_size_api.UpdateRingSize) (*response_model.Response, error) {
	if err := h.service.Update(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingSizeHandlers) Delete(ctx context.Context, size *ring_size_api.DeleteRingSize) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
