package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetAllStands(ctx context.Context, dto *pro_api.GetStandsRequest) (*pro_api.StandResponse, error) {
	stands, err := h.service.Stand.GetAll(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.StandResponse{Stands: stands}, nil
}

func (h *Handler) CreateStand(ctx context.Context, dto *pro_api.CreateStandRequest) (stand *pro_api.IdResponse, err error) {
	stand, err = h.service.Stand.Create(dto)
	if err != nil {
		return nil, err
	}

	return stand, nil
}

func (h *Handler) UpdateStand(ctx context.Context, dto *pro_api.UpdateStandRequest) (*pro_api.IdResponse, error) {
	err := h.service.Stand.Update(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteStand(ctx context.Context, dto *pro_api.DeleteStandRequest) (*pro_api.IdResponse, error) {
	err := h.service.Stand.Delete(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}
