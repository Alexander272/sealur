package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetPutgm(ctx context.Context, dto *proto.GetPutgmRequest) (*proto.PutgmResponse, error) {
	putgm, err := h.service.Putgm.Get(dto)
	if err != nil {
		return nil, err
	}

	return &proto.PutgmResponse{Putgm: putgm}, nil
}

func (h *Handler) CreatePutgm(ctx context.Context, dto *proto.CreatePutgmRequest) (*proto.IdResponse, error) {
	putg, err := h.service.Putgm.Create(dto)
	if err != nil {
		return nil, err
	}

	return putg, nil
}

func (h *Handler) UpdatePutgm(ctx context.Context, dto *proto.UpdatePutgmRequest) (*proto.IdResponse, error) {
	if err := h.service.Putgm.Update(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeletePutgm(ctx context.Context, dto *proto.DeletePutgmRequest) (*proto.IdResponse, error) {
	if err := h.service.Putgm.Delete(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}
