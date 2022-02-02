package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetAll(ctx context.Context, dto *proto.GetStandsRequest) (stands *proto.StandResponse, err error) {
	arr, err := h.service.GetAll(dto)
	if err != nil {
		return nil, err
	}

	stands = &proto.StandResponse{Stands: arr}

	return stands, nil
}

func (h *Handler) CreateStand(ctx context.Context, dto *proto.CreateStandRequest) (stand *proto.Id, err error) {
	stand, err = h.service.Create(dto)
	if err != nil {
		return nil, nil
	}

	return stand, nil
}

func (h *Handler) UpdateStand(ctx context.Context, dto *proto.UpdateStandRequest) (*proto.Id, error) {
	err := h.service.Update(dto)
	if err != nil {
		return nil, nil
	}

	return &proto.Id{Id: dto.Id}, nil
}

func (h *Handler) DeleteStand(ctx context.Context, dto *proto.DeleteStandRequest) (*proto.Id, error) {
	err := h.service.Delete(dto)
	if err != nil {
		return nil, nil
	}

	return &proto.Id{Id: dto.Id}, nil
}
