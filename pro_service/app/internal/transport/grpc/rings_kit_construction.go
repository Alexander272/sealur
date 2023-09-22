package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_construction_api"
)

type RingsKitConstructionHandlers struct {
	service service.RingsKitConstruction
	rings_kit_construction_api.UnimplementedRingsKitConstructionServiceServer
}

func NewRingsKitConstructionHandlers(service service.RingsKitConstruction) *RingsKitConstructionHandlers {
	return &RingsKitConstructionHandlers{
		service: service,
	}
}

func (h *RingsKitConstructionHandlers) GetAll(ctx context.Context, req *rings_kit_construction_api.GetRingsKitConstructions,
) (*rings_kit_construction_api.RingsKitConstructions, error) {
	constructions, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return &rings_kit_construction_api.RingsKitConstructions{RingsKitConstructions: constructions}, nil
}

func (h *RingsKitConstructionHandlers) Create(ctx context.Context, c *rings_kit_construction_api.CreateRingsKitConstruction) (*response_model.Response, error) {
	if err := h.service.Create(ctx, c); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingsKitConstructionHandlers) Update(ctx context.Context, m *rings_kit_construction_api.UpdateRingsKitConstruction) (*response_model.Response, error) {
	if err := h.service.Update(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingsKitConstructionHandlers) Delete(ctx context.Context, m *rings_kit_construction_api.DeleteRingsKitConstruction) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
