package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/rings_kit_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_api"
)

type RingsKitHandlers struct {
	service service.RingsKit
	rings_kit_api.UnimplementedRingsKitServiceServer
}

func NewRingsKitHandlers(service service.RingsKit) *RingsKitHandlers {
	return &RingsKitHandlers{
		service: service,
	}
}

func (h *RingsKitHandlers) Get(ctx context.Context, req *rings_kit_api.GetRingsKit) (*rings_kit_model.RingsKit, error) {
	ring, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return ring, nil
}
