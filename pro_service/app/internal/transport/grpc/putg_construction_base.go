package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_base_construction_api"
)

type PutgBaseConstructionHandlers struct {
	service service.PutgBaseConstruction
	putg_base_construction_api.UnimplementedPutgBaseConstructionServiceServer
}

func NewPutgBaseConstructionHandlers(service service.PutgBaseConstruction) *PutgBaseConstructionHandlers {
	return &PutgBaseConstructionHandlers{
		service: service,
	}
}

func (h *PutgBaseConstructionHandlers) Get(ctx context.Context, req *putg_base_construction_api.GetPutgBaseConstruction,
) (*putg_base_construction_api.PutgBaseConstruction, error) {
	constructions, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_base_construction_api.PutgBaseConstruction{Constructions: constructions}, nil
}

func (h *PutgBaseConstructionHandlers) Create(ctx context.Context, construction *putg_base_construction_api.CreatePutgBaseConstruction,
) (*response_model.Response, error) {
	if err := h.service.Create(ctx, construction); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgBaseConstructionHandlers) Update(ctx context.Context, construction *putg_base_construction_api.UpdatePutgBaseConstruction,
) (*response_model.Response, error) {
	if err := h.service.Update(ctx, construction); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgBaseConstructionHandlers) Delete(ctx context.Context, construction *putg_base_construction_api.DeletePutgBaseConstruction,
) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, construction); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
