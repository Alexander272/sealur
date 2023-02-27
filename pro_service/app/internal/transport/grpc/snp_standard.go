package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_standard_api"
)

type SnpStandardHandlers struct {
	service service.SnpStandard
	snp_standard_api.UnimplementedSnpStandardServiceServer
}

func NewSnpStandardHandlers(service service.SnpStandard) *SnpStandardHandlers {
	return &SnpStandardHandlers{
		service: service,
	}
}

// func (h *SnpStandardHandlers) Get(ctx context.Context, req *snp_standard_api.GetAllSnpStandards) (*snp_standard_api.SnpStandards, error) {
// 	standards, err := h.service.Get(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &snp_standard_api.SnpStandards{SnpStandards: standards}, nil
// }

func (h *SnpStandardHandlers) Create(ctx context.Context, standard *snp_standard_api.CreateSnpStandard) (*response_model.Response, error) {
	if err := h.service.Create(ctx, standard); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpStandardHandlers) CreateSeveral(ctx context.Context, standards *snp_standard_api.CreateSeveralSnpStandard) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, standards); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpStandardHandlers) Update(ctx context.Context, standard *snp_standard_api.UpdateSnpStandard) (*response_model.Response, error) {
	if err := h.service.Update(ctx, standard); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpStandardHandlers) Delete(ctx context.Context, standard *snp_standard_api.DeleteSnpStandard) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, standard); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
