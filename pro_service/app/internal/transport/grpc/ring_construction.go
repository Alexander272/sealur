package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_construction_api"
)

type RingConstructionHandlers struct {
	service service.RingConstruction
	*ring_construction_api.UnimplementedRingConstructionServiceServer
}

func NewRingConstructionHandlers(service service.RingConstruction) *RingConstructionHandlers {
	return &RingConstructionHandlers{
		service: service,
	}
}

func (h *RingConstructionHandlers) GetAll(ctx context.Context, req *ring_construction_api.GetRingConstructions,
) (*ring_construction_api.RingConstructions, error) {
	constructions, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ring_construction_api.RingConstructions{RingConstruction: constructions}, nil
}

func (h *RingConstructionHandlers) Create(ctx context.Context, m *ring_construction_api.CreateRingConstruction) (*response_model.Response, error) {
	if err := h.service.Create(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingConstructionHandlers) Update(ctx context.Context, m *ring_construction_api.UpdateRingConstruction) (*response_model.Response, error) {
	if err := h.service.Update(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingConstructionHandlers) Delete(ctx context.Context, m *ring_construction_api.DeleteRingConstruction) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
