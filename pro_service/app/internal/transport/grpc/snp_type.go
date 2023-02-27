package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_type_api"
)

type SnpTypeHandlers struct {
	service service.SnpType
	snp_type_api.UnimplementedSnpTypeServiceServer
}

func NewSnpTypeHandlers(service service.SnpType) *SnpTypeHandlers {
	return &SnpTypeHandlers{
		service: service,
	}
}

func (h *SnpTypeHandlers) Get(ctx context.Context, req *snp_type_api.GetSnpTypes) (*snp_type_api.SnpTypes, error) {
	snpTypes, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return &snp_type_api.SnpTypes{SnpTypes: snpTypes}, nil
}

func (h *SnpTypeHandlers) Create(ctx context.Context, snpType *snp_type_api.CreateSnpType) (*response_model.Response, error) {
	if err := h.service.Create(ctx, snpType); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpTypeHandlers) CreateSeveral(ctx context.Context, snpTypes *snp_type_api.CreateSeveralSnpType) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, snpTypes); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpTypeHandlers) Update(ctx context.Context, snpType *snp_type_api.UpdateSnpType) (*response_model.Response, error) {
	if err := h.service.Update(ctx, snpType); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpTypeHandlers) Delete(ctx context.Context, snpType *snp_type_api.DeleteSnpType) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, snpType); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
