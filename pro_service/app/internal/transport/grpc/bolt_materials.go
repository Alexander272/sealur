package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetBoltMaterials(ctx context.Context, req *pro_api.GetBoltMaterialsRequest) (*pro_api.BoltMaterialsResponse, error) {
	mats, err := h.service.BoltMaterials.GetAll(req)
	if err != nil {
		return nil, err
	}

	return &pro_api.BoltMaterialsResponse{Materials: mats}, nil
}

func (h *Handler) CreateBoltMaterials(ctx context.Context, mat *pro_api.CreateBoltMaterialsRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.BoltMaterials.Create(mat)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdateBoltMaterials(ctx context.Context, mat *pro_api.UpdateBoltMaterialsRequest) (*pro_api.IdResponse, error) {
	if err := h.service.BoltMaterials.Update(mat); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: mat.Id}, nil
}

func (h *Handler) DeleteBoltMaterials(ctx context.Context, mat *pro_api.DeleteBoltMaterialsRequest) (*pro_api.IdResponse, error) {
	if err := h.service.BoltMaterials.Delete(mat); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: mat.Id}, nil
}
