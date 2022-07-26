package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetSizes(ctx context.Context, dto *pro_api.GetSizesRequest) (*pro_api.SizeResponse, error) {
	sizes, dn, err := h.service.Size.Get(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.SizeResponse{Sizes: sizes, Dn: dn}, nil
}

func (h *Handler) GetAllSizes(ctx context.Context, dto *pro_api.GetSizesRequest) (*pro_api.SizeResponse, error) {
	sizes, dn, err := h.service.Size.GetAll(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.SizeResponse{Sizes: sizes, Dn: dn}, nil
}

func (h *Handler) CreateSize(ctx context.Context, dto *pro_api.CreateSizeRequest) (*pro_api.IdResponse, error) {
	size, err := h.service.Size.Create(dto)
	if err != nil {
		return nil, err
	}

	return size, nil
}

func (h *Handler) CreateManySizes(ctx context.Context, dto *pro_api.CreateSizesRequest) (*pro_api.IdResponse, error) {
	err := h.service.Size.CreateMany(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: ""}, nil
}

func (h *Handler) UpdateSize(ctx context.Context, dto *pro_api.UpdateSizeRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Size.Update(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteSize(ctx context.Context, dto *pro_api.DeleteSizeRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Size.Delete(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteAllSize(ctx context.Context, dto *pro_api.DeleteAllSizeRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Size.DeleteAll(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: ""}, nil
}
