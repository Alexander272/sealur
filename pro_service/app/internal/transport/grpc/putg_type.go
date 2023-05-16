package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_type_api"
)

type PutgTypeHandlers struct {
	service service.PutgType
	putg_type_api.UnimplementedPutgTypeServiceServer
}

func NewPutgTypeHandlers(service service.PutgType) *PutgTypeHandlers {
	return &PutgTypeHandlers{
		service: service,
	}
}

func (h *PutgTypeHandlers) Get(ctx context.Context, req *putg_type_api.GetPutgType) (*putg_type_api.PutgType, error) {
	types, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_type_api.PutgType{PutgTypes: types}, nil
}

func (h *PutgTypeHandlers) Create(ctx context.Context, putgType *putg_type_api.CreatePutgType) (*response_model.Response, error) {
	if err := h.service.Create(ctx, putgType); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgTypeHandlers) Update(ctx context.Context, putgType *putg_type_api.UpdatePutgType) (*response_model.Response, error) {
	if err := h.service.Update(ctx, putgType); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgTypeHandlers) Delete(ctx context.Context, putgType *putg_type_api.DeletePutgType) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, putgType); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
