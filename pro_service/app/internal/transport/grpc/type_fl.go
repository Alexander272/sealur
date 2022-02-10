package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetTypeFl(ctx context.Context, dto *proto.GetTypeFlRequest) (*proto.TypeFlResponse, error) {
	fl, err := h.service.TypeFl.Get()
	if err != nil {
		return nil, err
	}

	return &proto.TypeFlResponse{TypeFl: fl}, nil
}

func (h *Handler) CreateTypeFl(ctx context.Context, dto *proto.CreateTypeFlRequest) (*proto.IdResponse, error) {
	id, err := h.service.TypeFl.Create(dto)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdateTypeFl(ctx context.Context, dto *proto.UpdateTypeFlRequest) (*proto.IdResponse, error) {
	if err := h.service.TypeFl.Update(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteTypeFl(ctx context.Context, dto *proto.DeleteTypeFlRequest) (*proto.IdResponse, error) {
	if err := h.service.TypeFl.Delete(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}
