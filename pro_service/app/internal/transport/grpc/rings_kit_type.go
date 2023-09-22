package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_type_api"
)

type RingsKitTypeHandlers struct {
	service service.RingsKitType
	rings_kit_type_api.UnimplementedRingsKitTypeServiceServer
}

func NewRingsKitTypeHandlers(service service.RingsKitType) *RingsKitTypeHandlers {
	return &RingsKitTypeHandlers{
		service: service,
	}
}

func (h *RingsKitTypeHandlers) GetAll(ctx context.Context, req *rings_kit_type_api.GetRingsKitTypes) (*rings_kit_type_api.RingsKitTypes, error) {
	types, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return &rings_kit_type_api.RingsKitTypes{RingsKitTypes: types}, nil
}

func (h *RingsKitTypeHandlers) Create(ctx context.Context, kit *rings_kit_type_api.CreateRingsKitType) (*response_model.Response, error) {
	if err := h.service.Create(ctx, kit); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingsKitTypeHandlers) Update(ctx context.Context, m *rings_kit_type_api.UpdateRingsKitType) (*response_model.Response, error) {
	if err := h.service.Update(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *RingsKitTypeHandlers) Delete(ctx context.Context, m *rings_kit_type_api.DeleteRingsKitType) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, m); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
