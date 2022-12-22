package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetPutgmImage(ctx context.Context, dto *pro_api.GetPutgmImageRequest) (*pro_api.PutgmImageResponse, error) {
	images, err := h.service.PutgmImage.Get(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.PutgmImageResponse{PutgmImage: images}, nil
}

func (h *Handler) CreatePutgmImage(ctx context.Context, dto *pro_api.CreatePutgmImageRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.PutgmImage.Create(dto)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (h *Handler) UpdatePutgmImage(ctx context.Context, dto *pro_api.UpdatePutgmImageRequest) (*pro_api.IdResponse, error) {
	if err := h.service.PutgmImage.Update(dto); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeletePutgmImage(ctx context.Context, dto *pro_api.DeletePutgmImageRequest) (*pro_api.IdResponse, error) {
	if err := h.service.PutgmImage.Delete(dto); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}
