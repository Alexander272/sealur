package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetPutgImage(ctx context.Context, dto *pro_api.GetPutgImageRequest) (*pro_api.PutgImageResponse, error) {
	images, err := h.service.PutgImage.Get(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.PutgImageResponse{PutgImage: images}, nil
}

func (h *Handler) CreatePutgImage(ctx context.Context, dto *pro_api.CreatePutgImageRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.PutgImage.Create(dto)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdatePutgImage(ctx context.Context, dto *pro_api.UpdatePutgImageRequest) (*pro_api.IdResponse, error) {
	if err := h.service.PutgImage.Update(dto); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeletePutgImage(ctx context.Context, dto *pro_api.DeletePutgImageRequest) (*pro_api.IdResponse, error) {
	if err := h.service.PutgImage.Delete(dto); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}
