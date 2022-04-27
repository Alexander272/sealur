package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetPutgImage(ctx context.Context, dto *proto.GetPutgImageRequest) (*proto.PutgImageResponse, error) {
	images, err := h.service.PutgImage.Get(dto)
	if err != nil {
		return nil, err
	}

	return &proto.PutgImageResponse{PutgImage: images}, nil
}

func (h *Handler) CreatePutgImage(ctx context.Context, dto *proto.CreatePutgImageRequest) (*proto.IdResponse, error) {
	id, err := h.service.PutgImage.Create(dto)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdatePutgImage(ctx context.Context, dto *proto.UpdatePutgImageRequest) (*proto.IdResponse, error) {
	if err := h.service.PutgImage.Update(dto); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeletePutgImage(ctx context.Context, dto *proto.DeletePutgImageRequest) (*proto.IdResponse, error) {
	if err := h.service.PutgImage.Delete(dto); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}
