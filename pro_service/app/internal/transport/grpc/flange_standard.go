package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/flange_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
)

type FlangeStandardHandlers struct {
	service service.FlangeStandard
	flange_standard_api.UnimplementedFlangeStandardServiceServer
}

func NewFlangeStandardHandlers(service service.FlangeStandard) *FlangeStandardHandlers {
	return &FlangeStandardHandlers{
		service: service,
	}
}

func (h *FlangeStandardHandlers) GetAll(ctx context.Context, req *flange_standard_api.GetAllFlangeStandards) (*flange_standard_api.FlangeStandards, error) {
	flanges, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return &flange_standard_api.FlangeStandards{FlangeStandards: flanges}, nil
}

func (h *FlangeStandardHandlers) Create(ctx context.Context, flange *flange_standard_api.CreateFlangeStandard) (*response_model.Response, error) {
	if err := h.service.Create(ctx, flange); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *FlangeStandardHandlers) CreateSeveral(ctx context.Context, flanges *flange_standard_api.CreateSeveralFlangeStandard) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, flanges); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *FlangeStandardHandlers) Update(ctx context.Context, flange *flange_standard_api.UpdateFlangeStandard) (*response_model.Response, error) {
	if err := h.service.Update(ctx, flange); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *FlangeStandardHandlers) Delete(ctx context.Context, flange *flange_standard_api.DeleteFlangeStandard) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, flange); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
