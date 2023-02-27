package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/standard_api"
)

type StandardHandlers struct {
	service service.Standard
	standard_api.UnimplementedStandardServiceServer
}

func NewStandardHandlers(service service.Standard) *StandardHandlers {
	return &StandardHandlers{
		service: service,
	}
}

func (h *StandardHandlers) GetAll(ctx context.Context, req *standard_api.GetAllStandards) (*standard_api.Standards, error) {
	standards, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return &standard_api.Standards{Standards: standards}, nil
}

func (h *StandardHandlers) Create(ctx context.Context, standard *standard_api.CreateStandard) (*response_model.Response, error) {
	if err := h.service.Create(ctx, standard); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *StandardHandlers) CreateSeveral(ctx context.Context, standards *standard_api.CreateSeveralStandard) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, standards); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *StandardHandlers) Update(ctx context.Context, standard *standard_api.UpdateStandard) (*response_model.Response, error) {
	if err := h.service.Update(ctx, standard); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *StandardHandlers) Delete(ctx context.Context, standard *standard_api.DeleteStandard) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, standard); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
