package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_construction_api"
)

type PutgConstructionHandlers struct {
	service service.PutgConstruction
	putg_construction_api.UnimplementedPutgConstructionServiceServer
}

func NewPutgConstructionHandlers(service service.PutgConstruction) *PutgConstructionHandlers {
	return &PutgConstructionHandlers{
		service: service,
	}
}

func (h *PutgConstructionHandlers) Get(ctx context.Context, req *putg_construction_api.GetPutgConstruction) (*putg_construction_api.PutgConstruction, error) {
	constructions, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_construction_api.PutgConstruction{Constructions: constructions}, nil
}

func (h *PutgConstructionHandlers) Create(ctx context.Context, construction *putg_construction_api.CreatePutgConstruction) (*response_model.Response, error) {
	if err := h.service.Create(ctx, construction); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgConstructionHandlers) Update(ctx context.Context, construction *putg_construction_api.UpdatePutgConstruction) (*response_model.Response, error) {
	if err := h.service.Update(ctx, construction); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgConstructionHandlers) Delete(ctx context.Context, construction *putg_construction_api.DeletePutgConstruction) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, construction); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
