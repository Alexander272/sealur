package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetPutg(ctx context.Context, dto *proto.GetPutgRequest) (*proto.PutgResponse, error) {
	putg, err := h.service.Putg.Get(dto)
	if err != nil {
		return nil, err
	}

	return &proto.PutgResponse{Putg: putg}, nil
}

func (h *Handler) CreatePutg(ctx context.Context, dto *proto.CreatePutgRequest) (*proto.IdResponse, error) {
	putg, err := h.service.Putg.Create(dto)
	if err != nil {
		return nil, err
	}

	return putg, nil
}

func (h *Handler) UpdatePutg(ctx context.Context, dto *proto.UpdatePutgRequest) (*proto.IdResponse, error) {
	if err := h.service.Putg.Update(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeletePutg(ctx context.Context, dto *proto.DeletePutgRequest) (*proto.IdResponse, error) {
	if err := h.service.Putg.Delete(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}
