package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_api"
)

type PutgFillerHandlers struct {
	service service.PutgFiller
	putg_filler_api.UnimplementedPutgFillerServiceServer
}

func NewPutgFillerHandlers(service service.PutgFiller) *PutgFillerHandlers {
	return &PutgFillerHandlers{
		service: service,
	}
}

func (h *PutgFillerHandlers) Get(ctx context.Context, req *putg_filler_api.GetPutgFiller) (*putg_filler_api.PutgFiller, error) {
	fillers, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_filler_api.PutgFiller{Fillers: fillers}, nil
}

func (h *PutgFillerHandlers) Create(ctx context.Context, filler *putg_filler_api.CreatePutgFiller) (*response_model.Response, error) {
	if err := h.service.Create(ctx, filler); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgFillerHandlers) Update(ctx context.Context, filler *putg_filler_api.UpdatePutgFiller) (*response_model.Response, error) {
	if err := h.service.Update(ctx, filler); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgFillerHandlers) Delete(ctx context.Context, filler *putg_filler_api.DeletePutgFiller) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, filler); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
