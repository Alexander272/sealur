package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetBoltMaterials(ctx context.Context, req *proto.GetBoltMaterialsRequest) (*proto.BoltMaterialsResponse, error) {
	mats, err := h.service.BoltMaterials.GetAll(req)
	if err != nil {
		return nil, err
	}

	return &proto.BoltMaterialsResponse{Materials: mats}, nil
}

func (h *Handler) CreateBoltMaterials(ctx context.Context, mat *proto.CreateBoltMaterialsRequest) (*proto.IdResponse, error) {
	id, err := h.service.BoltMaterials.Create(mat)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdateBoltMaterials(ctx context.Context, mat *proto.UpdateBoltMaterialsRequest) (*proto.IdResponse, error) {
	if err := h.service.BoltMaterials.Update(mat); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: mat.Id}, nil
}

func (h *Handler) DeleteBoltMaterials(ctx context.Context, mat *proto.DeleteBoltMaterialsRequest) (*proto.IdResponse, error) {
	if err := h.service.BoltMaterials.Delete(mat); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: mat.Id}, nil
}
