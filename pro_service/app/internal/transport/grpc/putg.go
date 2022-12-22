package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetPutg(ctx context.Context, dto *pro_api.GetPutgRequest) (*pro_api.PutgResponse, error) {
	putg, err := h.service.Putg.Get(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.PutgResponse{Putg: putg}, nil
}

func (h *Handler) CreatePutg(ctx context.Context, dto *pro_api.CreatePutgRequest) (*pro_api.IdResponse, error) {
	putg, err := h.service.Putg.Create(dto)
	if err != nil {
		return nil, err
	}

	return putg, nil
}

func (h *Handler) UpdatePutg(ctx context.Context, dto *pro_api.UpdatePutgRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Putg.Update(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeletePutg(ctx context.Context, dto *pro_api.DeletePutgRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Putg.Delete(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}
