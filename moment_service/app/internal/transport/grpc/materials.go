package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment/material_api"
)

type MaterialsHandlers struct {
	service service.Materials
	material_api.UnimplementedMaterialsServiceServer
}

func NewMaterialsHandlers(service service.Materials) *MaterialsHandlers {
	return &MaterialsHandlers{
		service: service,
	}
}

func (h *MaterialsHandlers) GetMaterials(ctx context.Context, req *material_api.GetMaterialsRequest) (*material_api.MaterialsResponse, error) {
	materials, err := h.service.GetMaterials(ctx, req)
	if err != nil {
		return nil, err
	}

	return &material_api.MaterialsResponse{Materials: materials}, nil
}

func (h *MaterialsHandlers) GetMaterialsWithIsEmpty(ctx context.Context, req *material_api.GetMaterialsRequest,
) (*material_api.MaterialsWithIsEmptyResponse, error) {
	materials, err := h.service.GetMaterialsWithIsEmpty(ctx, req)
	if err != nil {
		return nil, err
	}

	return &material_api.MaterialsWithIsEmptyResponse{Materials: materials}, nil
}

func (h *MaterialsHandlers) GetMaterialsData(ctx context.Context, req *material_api.GetMaterialsDataRequest) (*material_api.MaterialsDataResponse, error) {
	materials, err := h.service.GetMaterialsData(ctx, req)
	if err != nil {
		return nil, err
	}

	return materials, nil
}

func (h *MaterialsHandlers) CreateMaterial(ctx context.Context, material *material_api.CreateMaterialRequest) (*material_api.IdResponse, error) {
	id, err := h.service.CreateMaterial(ctx, material)
	if err != nil {
		return nil, err
	}

	return &material_api.IdResponse{Id: id}, nil
}

func (h *MaterialsHandlers) UpdateMaterial(ctx context.Context, material *material_api.UpdateMaterialRequest) (*material_api.Response, error) {
	if err := h.service.UpdateMaterial(ctx, material); err != nil {
		return nil, err
	}
	return &material_api.Response{}, nil
}

func (h *MaterialsHandlers) DeleteMaterial(ctx context.Context, material *material_api.DeleteMaterialRequest) (*material_api.Response, error) {
	if err := h.service.DeleteMaterial(ctx, material); err != nil {
		return nil, err
	}
	return &material_api.Response{}, nil
}
