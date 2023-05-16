package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_material_api"
)

type PutgMaterialHandlers struct {
	service service.PutgMaterial
	putg_material_api.UnimplementedPutgMaterialServiceServer
}

func NewPutgMaterialHandlers(service service.PutgMaterial) *PutgMaterialHandlers {
	return &PutgMaterialHandlers{
		service: service,
	}
}

func (h *PutgMaterialHandlers) Get(ctx context.Context, req *putg_material_api.GetPutgMaterial) (*putg_material_api.PutgMaterials, error) {
	materials, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_material_api.PutgMaterials{Materials: materials}, nil
}

func (h *PutgMaterialHandlers) Create(ctx context.Context, material *putg_material_api.CreatePutgMaterial) (*response_model.Response, error) {
	if err := h.service.Create(ctx, material); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgMaterialHandlers) Update(ctx context.Context, material *putg_material_api.UpdatePutgMaterial) (*response_model.Response, error) {
	if err := h.service.Update(ctx, material); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgMaterialHandlers) Delete(ctx context.Context, material *putg_material_api.DeletePutgMaterial) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, material); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
