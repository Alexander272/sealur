package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetPutgmImage(ctx context.Context, dto *proto.GetPutgmImageRequest) (*proto.PutgmImageResponse, error) {
	images, err := h.service.PutgmImage.Get(dto)
	if err != nil {
		return nil, err
	}

	return &proto.PutgmImageResponse{PutgmImage: images}, nil
}

func (h *Handler) CreatePutgmImage(ctx context.Context, dto *proto.CreatePutgmImageRequest) (*proto.IdResponse, error) {
	id, err := h.service.PutgmImage.Create(dto)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdatePutgmImage(ctx context.Context, dto *proto.UpdatePutgmImageRequest) (*proto.IdResponse, error) {
	if err := h.service.PutgmImage.Update(dto); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeletePutgmImage(ctx context.Context, dto *proto.DeletePutgmImageRequest) (*proto.IdResponse, error) {
	if err := h.service.PutgmImage.Delete(dto); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}
