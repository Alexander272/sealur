package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_size_api"
)

type SnpSizeHandlers struct {
	service service.SnpSize
	snp_size_api.UnimplementedSnpSizeServiceServer
}

func NewSnpSizeHandlers(service service.SnpSize) *SnpSizeHandlers {
	return &SnpSizeHandlers{
		service: service,
	}
}

func (h *SnpSizeHandlers) Get(ctx context.Context, req *snp_size_api.GetSnpSize) (*snp_size_api.SnpSize, error) {
	sizes, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return &snp_size_api.SnpSize{Sizes: sizes}, nil
}

func (h *SnpSizeHandlers) Create(ctx context.Context, size *snp_size_api.CreateSnpSize) (*response_model.Response, error) {
	if err := h.service.Create(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpSizeHandlers) CreateSeveral(ctx context.Context, sizes *snp_size_api.CreateSeveralSnpSize) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, sizes); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpSizeHandlers) Update(ctx context.Context, size *snp_size_api.UpdateSnpSize) (*response_model.Response, error) {
	if err := h.service.Update(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpSizeHandlers) Delete(ctx context.Context, size *snp_size_api.DeleteSnpSize) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
