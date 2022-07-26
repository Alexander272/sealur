package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetMaterials(ctx context.Context, req *pro_api.GetMaterialsRequest) (*pro_api.MaterialsResponse, error) {
	mats, err := h.service.Materials.GetAll(req)
	if err != nil {
		return nil, err
	}

	return &pro_api.MaterialsResponse{Materials: mats}, nil
}

func (h *Handler) CreateMaterials(ctx context.Context, mat *pro_api.CreateMaterialsRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.Materials.Create(mat)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdateMaterials(ctx context.Context, mat *pro_api.UpdateMaterialsRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Materials.Update(mat); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: mat.Id}, nil
}

func (h *Handler) DeleteMaterials(ctx context.Context, mat *pro_api.DeleteMaterialsRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Materials.Delete(mat); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: mat.Id}, nil
}
