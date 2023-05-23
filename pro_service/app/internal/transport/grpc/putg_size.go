package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_size_api"
)

type PutgSizeHandlers struct {
	service service.PutgSize
	putg_size_api.UnimplementedPutgSizeServiceServer
}

func NewPutgSizeHandlers(service service.PutgSize) *PutgSizeHandlers {
	return &PutgSizeHandlers{
		service: service,
	}
}

func (h *PutgSizeHandlers) Get(ctx context.Context, req *putg_size_api.GetPutgSize) (*putg_size_api.PutgSize, error) {
	sizes, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_size_api.PutgSize{Sizes: sizes}, nil
}

func (h *PutgSizeHandlers) GetNew(ctx context.Context, req *putg_size_api.GetPutgSize_New) (*putg_size_api.PutgSize, error) {
	sizes, err := h.service.GetNew(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_size_api.PutgSize{Sizes: sizes}, nil
}

func (h *PutgSizeHandlers) Create(ctx context.Context, size *putg_size_api.CreatePutgSize) (*response_model.Response, error) {
	if err := h.service.Create(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgSizeHandlers) Update(ctx context.Context, size *putg_size_api.UpdatePutgSize) (*response_model.Response, error) {
	if err := h.service.Update(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgSizeHandlers) Delete(ctx context.Context, size *putg_size_api.DeletePutgSize) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, size); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
