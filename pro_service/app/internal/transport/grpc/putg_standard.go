package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_standard_api"
)

type PutgStandardHandlers struct {
	service service.PutgStandard
	putg_standard_api.UnimplementedPutgStandardServiceServer
}

func NewPutgStandardHandlers(service service.PutgStandard) *PutgStandardHandlers {
	return &PutgStandardHandlers{
		service: service,
	}
}

func (h *PutgStandardHandlers) Get(ctx context.Context, req *putg_standard_api.GetPutgStandard) (*putg_standard_api.PutgStandard, error) {
	standards, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_standard_api.PutgStandard{Standards: standards}, nil
}

func (h *PutgStandardHandlers) Create(ctx context.Context, standard *putg_standard_api.CreatePutgStandard) (*response_model.Response, error) {
	if err := h.service.Create(ctx, standard); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgStandardHandlers) Update(ctx context.Context, standard *putg_standard_api.UpdatePutgStandard) (*response_model.Response, error) {
	if err := h.service.Update(ctx, standard); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgStandardHandlers) Delete(ctx context.Context, standard *putg_standard_api.DeletePutgStandard) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, standard); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
