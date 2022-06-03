package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetMaterials(ctx context.Context, req *proto.GetMaterialsRequest) (*proto.MaterialsResponse, error) {
	mats, err := h.service.Materials.GetAll(req)
	if err != nil {
		return nil, err
	}

	return &proto.MaterialsResponse{Materials: mats}, nil
}

func (h *Handler) CreateMaterials(ctx context.Context, mat *proto.CreateMaterialsRequest) (*proto.IdResponse, error) {
	id, err := h.service.Materials.Create(mat)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdateMaterials(ctx context.Context, mat *proto.UpdateMaterialsRequest) (*proto.IdResponse, error) {
	if err := h.service.Materials.Update(mat); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: mat.Id}, nil
}

func (h *Handler) DeleteMaterials(ctx context.Context, mat *proto.DeleteMaterialsRequest) (*proto.IdResponse, error) {
	if err := h.service.Materials.Delete(mat); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: mat.Id}, nil
}
