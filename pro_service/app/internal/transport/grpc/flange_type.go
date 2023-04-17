package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
)

type FlangeTypeHandlers struct {
	service service.FlangeType
	flange_type_api.UnimplementedFlangeTypeServiceServer
}

func NewFlangeTypeHandlers(service service.FlangeType) *FlangeTypeHandlers {
	return &FlangeTypeHandlers{
		service: service,
	}
}

func (h *FlangeTypeHandlers) Get(ctx context.Context, req *flange_type_api.GetFlangeType) (*flange_type_api.FlangeType, error) {
	types, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return &flange_type_api.FlangeType{FlangeType: types}, nil
}

func (h *FlangeTypeHandlers) Create(ctx context.Context, flange *flange_type_api.CreateFlangeType) (*response_model.Response, error) {
	if err := h.service.Create(ctx, flange); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *FlangeTypeHandlers) Update(ctx context.Context, flange *flange_type_api.UpdateFlangeType) (*response_model.Response, error) {
	if err := h.service.Update(ctx, flange); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *FlangeTypeHandlers) Delete(ctx context.Context, flange *flange_type_api.DeleteFlangeType) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, flange); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
