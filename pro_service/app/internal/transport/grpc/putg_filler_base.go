package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_base_api"
)

type PutgBaseFillerHandlers struct {
	service service.PutgBaseFiller
	putg_filler_base_api.UnimplementedPutgBaseFillerServiceServer
}

func NewPutgBaseFillerHandlers(service service.PutgBaseFiller) *PutgBaseFillerHandlers {
	return &PutgBaseFillerHandlers{
		service: service,
	}
}

func (h *PutgBaseFillerHandlers) Get(ctx context.Context, req *putg_filler_base_api.GetPutgBaseFiller) (*putg_filler_base_api.PutgBaseFiller, error) {
	fillers, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_filler_base_api.PutgBaseFiller{Fillers: fillers}, nil
}

func (h *PutgBaseFillerHandlers) Create(ctx context.Context, filler *putg_filler_base_api.CreatePutgBaseFiller) (*response_model.Response, error) {
	if err := h.service.Create(ctx, filler); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgBaseFillerHandlers) Update(ctx context.Context, filler *putg_filler_base_api.UpdatePutgBaseFiller) (*response_model.Response, error) {
	if err := h.service.Update(ctx, filler); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgBaseFillerHandlers) Delete(ctx context.Context, filler *putg_filler_base_api.DeletePutgBaseFiller) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, filler); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
