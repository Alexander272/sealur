package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_modifying_api"
)

type RingModifyingHandlers struct {
	service service.RingModifying
	ring_modifying_api.UnimplementedRingModifyingServiceServer
}

func NewRingModifyingHandlers(service service.RingModifying) *RingModifyingHandlers {
	return &RingModifyingHandlers{
		service: service,
	}
}

func (h *RingModifyingHandlers) GetAll(ctx context.Context, req *ring_modifying_api.GetRingModifying) (*ring_modifying_api.RingModifying, error) {
	mods, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ring_modifying_api.RingModifying{Modifying: mods}, nil
}

func (h *RingModifyingHandlers) Create(ctx context.Context, m *ring_modifying_api.CreateRingModifying) (*response_model.Response, error) {
	if err := h.service.Create(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingModifyingHandlers) Update(ctx context.Context, m *ring_modifying_api.UpdateRingModifying) (*response_model.Response, error) {
	if err := h.service.Update(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingModifyingHandlers) Delete(ctx context.Context, m *ring_modifying_api.DeleteRingModifying) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
