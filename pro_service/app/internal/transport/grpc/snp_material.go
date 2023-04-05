package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
)

type SnpMaterialHandlers struct {
	service service.SnpMaterial
	snp_material_api.UnimplementedSnpMaterialServiceServer
}

func NewSnpMaterialHandlers(service service.SnpMaterial) *SnpMaterialHandlers {
	return &SnpMaterialHandlers{
		service: service,
	}
}

// func (h *SnpMaterialHandlers) Get(ctx context.Context, req *snp_material_api.GetSnpMaterial) (*snp_material_api.SnpMaterials, error) {
// 	materials, err := h.service.Get(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &snp_material_api.SnpMaterials{Materials: materials}, nil
// }

func (h *SnpMaterialHandlers) Create(ctx context.Context, material *snp_material_api.CreateSnpMaterial) (*response_model.Response, error) {
	if err := h.service.Create(ctx, material); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpMaterialHandlers) Update(ctx context.Context, material *snp_material_api.UpdateSnpMaterial) (*response_model.Response, error) {
	if err := h.service.Update(ctx, material); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpMaterialHandlers) Delete(ctx context.Context, material *snp_material_api.DeleteSnpMaterial) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, material); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
