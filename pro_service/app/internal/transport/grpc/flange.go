package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetAllFlanges(ctx context.Context, dto *proto.GetAllFlangeRequest) (*proto.FlangeResponse, error) {
	flanges, err := h.service.Flange.GetAll()
	if err != nil {
		return nil, err
	}

	return &proto.FlangeResponse{Flanges: flanges}, nil
}

func (h *Handler) CreateFlange(ctx context.Context, dto *proto.CreateFlangeRequest) (*proto.IdResponse, error) {
	flange, err := h.service.Flange.Create(dto)
	if err != nil {
		return nil, err
	}

	return flange, nil
}

func (h *Handler) UpdateFlange(ctx context.Context, dto *proto.UpdateFlangeRequest) (*proto.IdResponse, error) {
	if err := h.service.Flange.Update(dto); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteFlange(ctx context.Context, dto *proto.DeleteFlangeRequest) (*proto.IdResponse, error) {
	if err := h.service.Flange.Delete(dto); err != nil {
		return nil, err
	}

	return &proto.IdResponse{Id: dto.Id}, nil
}
