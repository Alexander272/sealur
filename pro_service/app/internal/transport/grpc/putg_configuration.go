package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_conf_api"
)

type PutgConfigurationHandlers struct {
	service service.PutgConfiguration
	putg_conf_api.UnimplementedPutgConfigurationServiceServer
}

func NewPutgConfigurationHandlers(service service.PutgConfiguration) *PutgConfigurationHandlers {
	return &PutgConfigurationHandlers{
		service: service,
	}
}

func (h *PutgConfigurationHandlers) Get(ctx context.Context, req *putg_conf_api.GetPutgConfiguration) (*putg_conf_api.PutgConfiguration, error) {
	configurations, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &putg_conf_api.PutgConfiguration{Configurations: configurations}, nil
}

func (h *PutgConfigurationHandlers) Create(ctx context.Context, conf *putg_conf_api.CreatePutgConfiguration) (*response_model.Response, error) {
	if err := h.service.Create(ctx, conf); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgConfigurationHandlers) Update(ctx context.Context, conf *putg_conf_api.UpdatePutgConfiguration) (*response_model.Response, error) {
	if err := h.service.Update(ctx, conf); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PutgConfigurationHandlers) Delete(ctx context.Context, conf *putg_conf_api.DeletePutgConfiguration) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, conf); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
