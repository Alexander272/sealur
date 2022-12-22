package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetAllFlanges(ctx context.Context, dto *pro_api.GetAllFlangeRequest) (*pro_api.FlangeResponse, error) {
	flanges, err := h.service.Flange.GetAll()
	if err != nil {
		return nil, err
	}

	return &pro_api.FlangeResponse{Flanges: flanges}, nil
}

func (h *Handler) CreateFlange(ctx context.Context, dto *pro_api.CreateFlangeRequest) (*pro_api.IdResponse, error) {
	flange, err := h.service.Flange.Create(dto)
	if err != nil {
		return nil, err
	}

	return flange, nil
}

func (h *Handler) UpdateFlange(ctx context.Context, dto *pro_api.UpdateFlangeRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Flange.Update(dto); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteFlange(ctx context.Context, dto *pro_api.DeleteFlangeRequest) (*pro_api.IdResponse, error) {
	if err := h.service.Flange.Delete(dto); err != nil {
		return nil, err
	}

	return &pro_api.IdResponse{Id: dto.Id}, nil
}
