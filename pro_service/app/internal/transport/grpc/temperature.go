package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/temperature_api"
)

type TemperatureHandlers struct {
	service service.Temperature
	temperature_api.UnimplementedTemperatureServiceServer
}

func NewTemperatureHandlers(service service.Temperature) *TemperatureHandlers {
	return &TemperatureHandlers{
		service: service,
	}
}

func (h *TemperatureHandlers) GetAll(ctx context.Context, req *temperature_api.GetAllTemperatures) (*temperature_api.Temperatures, error) {
	temperatures, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return &temperature_api.Temperatures{Temperatures: temperatures}, nil
}

func (h *TemperatureHandlers) Create(ctx context.Context, temperature *temperature_api.CreateTemperature) (*response_model.Response, error) {
	if err := h.service.Create(ctx, temperature); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *TemperatureHandlers) CreateSeveral(ctx context.Context, temperatures *temperature_api.CreateSeveralTemperature) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, temperatures); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *TemperatureHandlers) Update(ctx context.Context, temperature *temperature_api.UpdateTemperature) (*response_model.Response, error) {
	if err := h.service.Update(ctx, temperature); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *TemperatureHandlers) Delete(ctx context.Context, temperature *temperature_api.DeleteTemperature) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, temperature); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
