package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_density_api"
)

type RingDensityHandlers struct {
	service service.RingDensity
	*ring_density_api.UnimplementedRingDensityServiceServer
}

func NewRingDensityHandlers(service service.RingDensity) *RingDensityHandlers {
	return &RingDensityHandlers{
		service: service,
	}
}

func (h *RingDensityHandlers) GetAll(ctx context.Context, req *ring_density_api.GetRingDensity) (*ring_density_api.RingDensity, error) {
	density, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ring_density_api.RingDensity{Density: density}, nil
}

func (h *RingDensityHandlers) Create(ctx context.Context, m *ring_density_api.CreateRingDensity) (*response_model.Response, error) {
	if err := h.service.Create(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingDensityHandlers) Update(ctx context.Context, m *ring_density_api.UpdateRingDensity) (*response_model.Response, error) {
	if err := h.service.Update(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingDensityHandlers) Delete(ctx context.Context, m *ring_density_api.DeleteRingDensity) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
