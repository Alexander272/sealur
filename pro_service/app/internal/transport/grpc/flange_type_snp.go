package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
)

type FlangeTypeSnpHandlers struct {
	service service.FlangeTypeSnp
	flange_type_snp_api.UnimplementedFlangeTypeSnpServiceServer
}

func NewFlangeTypeSnpHandlers(service service.FlangeTypeSnp) *FlangeTypeSnpHandlers {
	return &FlangeTypeSnpHandlers{
		service: service,
	}
}

func (h *FlangeTypeSnpHandlers) Get(ctx context.Context, req *flange_type_snp_api.GetFlangeTypeSnp) (*flange_type_snp_api.FlangeType, error) {
	types, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return &flange_type_snp_api.FlangeType{FlangeTypeSnp: types}, nil
}

func (h *FlangeTypeSnpHandlers) Create(ctx context.Context, flange *flange_type_snp_api.CreateFlangeTypeSnp) (*response_model.Response, error) {
	if err := h.service.Create(ctx, flange); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *FlangeTypeSnpHandlers) CreateSeveral(ctx context.Context, flanges *flange_type_snp_api.CreateSeveralFlangeTypeSnp) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, flanges); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *FlangeTypeSnpHandlers) Update(ctx context.Context, flange *flange_type_snp_api.UpdateFlangeTypeSnp) (*response_model.Response, error) {
	if err := h.service.Update(ctx, flange); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *FlangeTypeSnpHandlers) Delete(ctx context.Context, flange *flange_type_snp_api.DeleteFlangeTypeSnp) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, flange); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
