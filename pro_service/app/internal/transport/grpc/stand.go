package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetAllStands(ctx context.Context, dto *proto.GetStandsRequest) (*proto.StandResponse, error) {
	stands, err := h.service.Stand.GetAll(dto)
	if err != nil {
		return nil, err
	}

	return &proto.StandResponse{Stands: stands}, nil
}

func (h *Handler) CreateStand(ctx context.Context, dto *proto.CreateStandRequest) (stand *proto.IdResponse, err error) {
	stand, err = h.service.Stand.Create(dto)
	if err != nil {
		return nil, err
	}

	return stand, nil
}

func (h *Handler) UpdateStand(ctx context.Context, dto *proto.UpdateStandRequest) (*proto.IdResponse, error) {
	err := h.service.Stand.Update(dto)
	if err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteStand(ctx context.Context, dto *proto.DeleteStandRequest) (*proto.IdResponse, error) {
	err := h.service.Stand.Delete(dto)
	if err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}
