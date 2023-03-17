package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
)

type PositionHandlers struct {
	service service.Position
	position_api.UnimplementedPositionServiceServer
}

func NewPositionHandlers(service service.Position) *PositionHandlers {
	return &PositionHandlers{
		service: service,
	}
}

func (h *PositionHandlers) Create(ctx context.Context, position *position_api.CreatePosition) (*response_model.IdResponse, error) {
	id, err := h.service.Create(ctx, position.Position)
	if err != nil {
		return nil, err
	}
	return &response_model.IdResponse{Id: id}, nil
}

func (h *PositionHandlers) Update(ctx context.Context, position *position_api.UpdatePosition) (*response_model.Response, error) {
	if err := h.service.Update(ctx, position.Position); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *PositionHandlers) Delete(ctx context.Context, position *position_api.DeletePosition) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, position.Id); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
