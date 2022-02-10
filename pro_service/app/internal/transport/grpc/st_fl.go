package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetStFl(ctx context.Context, dto *proto.GetStFlRequest) (*proto.StFlResponse, error) {
	st, err := h.service.StFl.Get()
	if err != nil {
		return nil, err
	}

	return &proto.StFlResponse{Stfl: st}, nil
}

func (h *Handler) CreateStFl(ctx context.Context, dto *proto.CreateStFlRequest) (*proto.IdResponse, error) {
	id, err := h.service.StFl.Create(dto)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdateStFl(ctx context.Context, dto *proto.UpdateStFlRequest) (*proto.IdResponse, error) {
	if err := h.service.StFl.Update(dto); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteStFl(ctx context.Context, dto *proto.DeleteStFlRequest) (*proto.IdResponse, error) {
	if err := h.service.StFl.Delete(dto); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}
