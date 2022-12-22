package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetSizeInt(ctx context.Context, req *pro_api.GetSizesIntRequest) (*pro_api.SizeIntResponse, error) {
	sizes, dn, err := h.service.SizeInt.Get(req)
	if err != nil {
		return nil, err
	}

	return &pro_api.SizeIntResponse{Sizes: sizes, Dn: dn}, nil
}

func (h *Handler) GetAllSizeInt(ctx context.Context, req *pro_api.GetAllSizeIntRequest) (*pro_api.SizeIntResponse, error) {
	sizes, dn, err := h.service.SizeInt.GetAll(req)
	if err != nil {
		return nil, err
	}

	return &pro_api.SizeIntResponse{Sizes: sizes, Dn: dn}, nil
}

func (h *Handler) CreateSizeInt(ctx context.Context, size *pro_api.CreateSizeIntRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.SizeInt.Create(size)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (h *Handler) CreateManySizesInt(ctx context.Context, dto *pro_api.CreateSizesIntRequest) (*pro_api.IdResponse, error) {
	err := h.service.SizeInt.CreateMany(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: ""}, nil
}

func (h *Handler) UpdateSizeInt(ctx context.Context, size *pro_api.UpdateSizeIntRequest) (*pro_api.IdResponse, error) {
	if err := h.service.SizeInt.Update(size); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: size.Id}, nil
}

func (h *Handler) DeleteSizeInt(ctx context.Context, size *pro_api.DeleteSizeIntRequest) (*pro_api.IdResponse, error) {
	if err := h.service.SizeInt.Delete(size); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: size.Id}, nil
}

func (h *Handler) DeleteAllSizeInt(ctx context.Context, size *pro_api.DeleteAllSizeIntRequest) (*pro_api.SuccessResponse, error) {
	if err := h.service.SizeInt.DeleteAll(size); err != nil {
		return nil, err
	}

	return &pro_api.SuccessResponse{Success: true}, nil
}
