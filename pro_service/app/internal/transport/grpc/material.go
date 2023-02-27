package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
)

type MaterialHandlers struct {
	service service.Material
	material_api.UnimplementedMaterialServiceServer
}

func NewMaterialHandlers(service service.Material) *MaterialHandlers {
	return &MaterialHandlers{
		service: service,
	}
}

func (h *MaterialHandlers) GetAll(ctx context.Context, req *material_api.GetAllMaterials) (*material_api.Materials, error) {
	materials, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return &material_api.Materials{Materials: materials}, nil
}

func (h *MaterialHandlers) Create(ctx context.Context, material *material_api.CreateMaterial) (*response_model.Response, error) {
	if err := h.service.Create(ctx, material); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *MaterialHandlers) CreateSeveral(ctx context.Context, materials *material_api.CreateSeveralMaterial) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, materials); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *MaterialHandlers) Update(ctx context.Context, material *material_api.UpdateMaterial) (*response_model.Response, error) {
	if err := h.service.Update(ctx, material); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *MaterialHandlers) Delete(ctx context.Context, material *material_api.DeleteMaterial) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, material); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
