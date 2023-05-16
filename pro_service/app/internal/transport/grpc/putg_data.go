package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_data_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_data_api"
)

type PutgDataHandlers struct {
	service service.PutgData
	putg_data_api.UnimplementedPutgDataServiceServer
}

func NewPutgDataHandlers(service service.PutgData) *PutgDataHandlers {
	return &PutgDataHandlers{
		service: service,
	}
}

func (h *PutgDataHandlers) Get(ctx context.Context, req *putg_data_api.GetPutgData) (*putg_data_api.PutgData, error) {
	data, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	//TODO тут не массив должен возвращаться
	return &putg_data_api.PutgData{Data: []*putg_data_model.PutgData{data}}, nil
}

func (h *PutgDataHandlers) Create(ctx context.Context, data *putg_data_api.CreatePutgData) (*response_model.Response, error) {
	if err := h.service.Create(ctx, data); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgDataHandlers) Update(ctx context.Context, data *putg_data_api.UpdatePutgData) (*response_model.Response, error) {
	if err := h.service.Update(ctx, data); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgDataHandlers) Delete(ctx context.Context, data *putg_data_api.DeletePutgData) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, data); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
