package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetSizeInt(ctx context.Context, req *proto.GetSizesIntRequest) (*proto.SizeIntResponse, error) {
	sizes, err := h.service.SizeInt.Get(req)
	if err != nil {
		return nil, err
	}

	return &proto.SizeIntResponse{Sizes: sizes}, nil
}

func (h *Handler) CreateSizeInt(ctx context.Context, size *proto.CreateSizeIntRequest) (*proto.IdResponse, error) {
	id, err := h.service.SizeInt.Create(size)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (h *Handler) UpdateSizeInt(ctx context.Context, size *proto.UpdateSizeIntRequest) (*proto.IdResponse, error) {
	if err := h.service.SizeInt.Update(size); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: size.Id}, nil
}

func (h *Handler) DeleteSizeInt(ctx context.Context, size *proto.DeleteSizeIntRequest) (*proto.IdResponse, error) {
	if err := h.service.SizeInt.Delete(size); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: size.Id}, nil
}
