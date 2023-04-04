package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
)

type SnpFillerHandlers struct {
	service service.SnpFiller
	snp_filler_api.UnimplementedSnpFillerServiceServer
}

func NewSnpFillerHandlers(service service.SnpFiller) *SnpFillerHandlers {
	return &SnpFillerHandlers{
		service: service,
	}
}

//TODO надо возвращать новые данные
// func (h *SnpFillerHandlers) Get(ctx context.Context, req *snp_filler_api.GetSnpFillers) (*snp_filler_api.SnpFillers, error) {
// 	fillers, err := h.service.GetAll(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &snp_filler_api.SnpFillers{SnpFillers: fillers}, nil
// }

func (h *SnpFillerHandlers) Create(ctx context.Context, filler *snp_filler_api.CreateSnpFiller) (*response_model.Response, error) {
	if err := h.service.Create(ctx, filler); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpFillerHandlers) CreateSeveral(ctx context.Context, fillers *snp_filler_api.CreateSeveralSnpFiller) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, fillers); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpFillerHandlers) Update(ctx context.Context, filler *snp_filler_api.UpdateSnpFiller) (*response_model.Response, error) {
	if err := h.service.Update(ctx, filler); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpFillerHandlers) Delete(ctx context.Context, filler *snp_filler_api.DeleteSnpFiller) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, filler); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
