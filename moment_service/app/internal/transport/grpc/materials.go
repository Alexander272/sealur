package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type MaterialsHandlers struct {
	service service.Materials
}

func NewMaterialsHandlers(service service.Materials) *MaterialsHandlers {
	return &MaterialsHandlers{
		service: service,
	}
}

func (h *MaterialsHandlers) GetMaterials(ctx context.Context, req *moment_proto.GetMaterialsRequest) (*moment_proto.MaterialsResponse, error) {
	materials, err := h.service.GetMaterials(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.MaterialsResponse{Materials: materials}, nil
}

func (h *MaterialsHandlers) GetMaterialsWithIsEmpty(ctx context.Context, req *moment_proto.GetMaterialsRequest,
) (*moment_proto.MaterialsWithIsEmptyResponse, error) {
	materials, err := h.service.GetMaterialsWithIsEmpty(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.MaterialsWithIsEmptyResponse{Materials: materials}, nil
}

func (h *MaterialsHandlers) GetMaterialsData(ctx context.Context, req *moment_proto.GetMaterialsDataRequest) (*moment_proto.MaterialsDataResponse, error) {
	materials, err := h.service.GetMaterialsData(ctx, req)
	if err != nil {
		return nil, err
	}

	return materials, nil
}

func (h *MaterialsHandlers) CreateMaterial(ctx context.Context, material *moment_proto.CreateMaterialRequest) (*moment_proto.IdResponse, error) {
	id, err := h.service.CreateMaterial(ctx, material)
	if err != nil {
		return nil, err
	}

	return &moment_proto.IdResponse{Id: id}, nil
}

func (h *MaterialsHandlers) UpdateMaterial(ctx context.Context, material *moment_proto.UpdateMaterialRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateMaterial(ctx, material); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *MaterialsHandlers) DeleteMaterial(ctx context.Context, material *moment_proto.DeleteMaterialRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteMaterial(ctx, material); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
