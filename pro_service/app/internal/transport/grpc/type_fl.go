package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetTypeFl(ctx context.Context, dto *pro_api.GetTypeFlRequest) (*pro_api.TypeFlResponse, error) {
	fl, err := h.service.TypeFl.Get()
	if err != nil {
		return nil, err
	}

	return &pro_api.TypeFlResponse{TypeFl: fl}, nil
}

func (h *Handler) GetAllTypeFl(ctx context.Context, dto *pro_api.GetTypeFlRequest) (*pro_api.TypeFlResponse, error) {
	fl, err := h.service.TypeFl.GetAll()
	if err != nil {
		return nil, err
	}

	return &pro_api.TypeFlResponse{TypeFl: fl}, nil
}

func (h *Handler) CreateTypeFl(ctx context.Context, dto *pro_api.CreateTypeFlRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.TypeFl.Create(dto)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdateTypeFl(ctx context.Context, dto *pro_api.UpdateTypeFlRequest) (*pro_api.IdResponse, error) {
	if err := h.service.TypeFl.Update(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteTypeFl(ctx context.Context, dto *pro_api.DeleteTypeFlRequest) (*pro_api.IdResponse, error) {
	if err := h.service.TypeFl.Delete(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}
