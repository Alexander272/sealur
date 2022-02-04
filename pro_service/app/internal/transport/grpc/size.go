package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetSizes(ctx context.Context, dto *proto.GetSizesRequest) (*proto.SizeResponse, error) {
	sizes, err := h.service.Size.Get(dto)
	if err != nil {
		return nil, err
	}

	return &proto.SizeResponse{Sizes: sizes}, nil
}

func (h *Handler) CreateSize(ctx context.Context, dto *proto.CreateSizeRequest) (*proto.IdResponse, error) {
	size, err := h.service.Size.Create(dto)
	if err != nil {
		return nil, err
	}

	return size, nil
}

func (h *Handler) UpdateSize(ctx context.Context, dto *proto.UpdateSizeRequest) (*proto.IdResponse, error) {
	if err := h.service.Size.Update(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteSize(ctx context.Context, dto *proto.DeleteSizeRequest) (*proto.IdResponse, error) {
	if err := h.service.Size.Delete(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}
