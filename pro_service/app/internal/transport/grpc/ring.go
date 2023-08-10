package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_api"
)

type RingHandlers struct {
	service service.Ring
	ring_api.UnimplementedRingServiceServer
}

func NewRingHandlers(service service.Ring) *RingHandlers {
	return &RingHandlers{
		service: service,
	}
}

func (h *RingHandlers) Get(ctx context.Context, req *ring_api.GetRings) (*ring_model.Ring, error) {
	ring, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return ring, nil
}
