package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetStFl(ctx context.Context, dto *pro_api.GetStFlRequest) (*pro_api.StFlResponse, error) {
	st, err := h.service.StFl.Get()
	if err != nil {
		return nil, err
	}

	return &pro_api.StFlResponse{Stfl: st}, nil
}

func (h *Handler) CreateStFl(ctx context.Context, dto *pro_api.CreateStFlRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.StFl.Create(dto)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdateStFl(ctx context.Context, dto *pro_api.UpdateStFlRequest) (*pro_api.IdResponse, error) {
	if err := h.service.StFl.Update(dto); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteStFl(ctx context.Context, dto *pro_api.DeleteStFlRequest) (*pro_api.IdResponse, error) {
	if err := h.service.StFl.Delete(dto); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}
