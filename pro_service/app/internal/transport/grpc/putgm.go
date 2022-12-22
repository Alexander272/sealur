package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetPutgm(ctx context.Context, dto *pro_api.GetPutgmRequest) (*pro_api.PutgmResponse, error) {
	putgm, err := h.service.Putgm.Get(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.PutgmResponse{Putgm: putgm}, nil
}

func (h *Handler) CreatePutgm(ctx context.Context, dto *pro_api.CreatePutgmRequest) (*pro_api.IdResponse, error) {
	putg, err := h.service.Putgm.Create(dto)
	if err != nil {
		return nil, err
	}

	return putg, nil
}

func (h *Handler) UpdatePutgm(ctx context.Context, dto *pro_api.UpdatePutgmRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Putgm.Update(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeletePutgm(ctx context.Context, dto *pro_api.DeletePutgmRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Putgm.Delete(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}
