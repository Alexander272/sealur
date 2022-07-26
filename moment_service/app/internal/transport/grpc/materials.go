package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type MaterialsHandlers struct {
	service service.Materials
	moment_api.UnimplementedMaterialsServiceServer
}

func NewMaterialsHandlers(service service.Materials) *MaterialsHandlers {
	return &MaterialsHandlers{
		service: service,
	}
}

func (h *MaterialsHandlers) GetMaterials(ctx context.Context, req *moment_api.GetMaterialsRequest) (*moment_api.MaterialsResponse, error) {
	materials, err := h.service.GetMaterials(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_api.MaterialsResponse{Materials: materials}, nil
}

func (h *MaterialsHandlers) GetMaterialsWithIsEmpty(ctx context.Context, req *moment_api.GetMaterialsRequest,
) (*moment_api.MaterialsWithIsEmptyResponse, error) {
	materials, err := h.service.GetMaterialsWithIsEmpty(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_api.MaterialsWithIsEmptyResponse{Materials: materials}, nil
}

func (h *MaterialsHandlers) GetMaterialsData(ctx context.Context, req *moment_api.GetMaterialsDataRequest) (*moment_api.MaterialsDataResponse, error) {
	materials, err := h.service.GetMaterialsData(ctx, req)
	if err != nil {
		return nil, err
	}

	return materials, nil
}

func (h *MaterialsHandlers) CreateMaterial(ctx context.Context, material *moment_api.CreateMaterialRequest) (*moment_api.IdResponse, error) {
	id, err := h.service.CreateMaterial(ctx, material)
	if err != nil {
		return nil, err
	}

	return &moment_api.IdResponse{Id: id}, nil
}

func (h *MaterialsHandlers) UpdateMaterial(ctx context.Context, material *moment_api.UpdateMaterialRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateMaterial(ctx, material); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *MaterialsHandlers) DeleteMaterial(ctx context.Context, material *moment_api.DeleteMaterialRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteMaterial(ctx, material); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}
