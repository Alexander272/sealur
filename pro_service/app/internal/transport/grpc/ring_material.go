package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_material_api"
)

type RingMaterialHandlers struct {
	service service.RingMaterial
	*ring_material_api.UnimplementedRingMaterialServiceServer
}

func NewRingMaterialHandlers(service service.RingMaterial) *RingMaterialHandlers {
	return &RingMaterialHandlers{
		service: service,
	}
}

func (h *RingMaterialHandlers) Get(ctx context.Context, req *ring_material_api.GetRingMaterial) (*ring_material_api.RingMaterial, error) {
	materials, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ring_material_api.RingMaterial{Materials: materials}, nil
}

func (h *RingMaterialHandlers) Create(ctx context.Context, m *ring_material_api.CreateRingMaterial) (*response_model.Response, error) {
	if err := h.service.Create(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingMaterialHandlers) Update(ctx context.Context, m *ring_material_api.UpdateRingMaterial) (*response_model.Response, error) {
	if err := h.service.Update(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingMaterialHandlers) Delete(ctx context.Context, m *ring_material_api.DeleteRingMaterial) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
