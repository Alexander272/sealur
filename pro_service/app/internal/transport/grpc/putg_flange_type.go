package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_flange_type_api"
)

type PutgFlangeTypeHandlers struct {
	service service.PutgFlangeType
	putg_flange_type_api.UnimplementedPutgFlangeTypeServiceServer
}

func NewPutgFlangeTypeHandlers(service service.PutgFlangeType) *PutgFlangeTypeHandlers {
	return &PutgFlangeTypeHandlers{
		service: service,
	}
}

func (h *PutgFlangeTypeHandlers) Get(ctx context.Context, req *putg_flange_type_api.GetPutgFlangeType) (*putg_flange_type_api.PutgFlangeType, error) {
	flangeTypes, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_flange_type_api.PutgFlangeType{FlangeTypes: flangeTypes}, nil
}

func (h *PutgFlangeTypeHandlers) Create(ctx context.Context, flangeType *putg_flange_type_api.CreatePutgFlangeType) (*response_model.Response, error) {
	if err := h.service.Create(ctx, flangeType); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgFlangeTypeHandlers) Update(ctx context.Context, flangeType *putg_flange_type_api.UpdatePutgFlangeType) (*response_model.Response, error) {
	if err := h.service.Update(ctx, flangeType); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgFlangeTypeHandlers) Delete(ctx context.Context, flangeType *putg_flange_type_api.DeletePutgFlangeType) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, flangeType); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
