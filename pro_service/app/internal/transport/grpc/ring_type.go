package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_type_api"
)

type RingTypeHandlers struct {
	service service.RingType
	*ring_type_api.UnimplementedSnpMaterialServiceServer
}

func NewRingTypeHandlers(service service.RingType) *RingTypeHandlers {
	return &RingTypeHandlers{
		service: service,
	}
}

func (h *RingTypeHandlers) GetAll(ctx context.Context, req *ring_type_api.GetRingTypes) (*ring_type_api.RingTypes, error) {
	ringTypes, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ring_type_api.RingTypes{RingTypes: ringTypes}, nil
}

func (h *RingTypeHandlers) Create(ctx context.Context, c *ring_type_api.CreateRingType) (*response_model.Response, error) {
	if err := h.service.Create(ctx, c); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingTypeHandlers) Update(ctx context.Context, c *ring_type_api.UpdateRingType) (*response_model.Response, error) {
	if err := h.service.Update(ctx, c); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingTypeHandlers) Delete(ctx context.Context, c *ring_type_api.DeleteRingType) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, c); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
