package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
)

func (h *FlangeHandlers) GetTypeFlange(ctx context.Context, req *flange_api.GetTypeFlangeRequest) (*flange_api.TypeFlangeResponse, error) {
	typeFlange, err := h.service.GetTypeFlange(ctx, req)
	if err != nil {
		return nil, err
	}

	return &flange_api.TypeFlangeResponse{TypeFlanges: typeFlange}, nil
}

func (h *FlangeHandlers) CreateTypeFlange(ctx context.Context, typeFlange *flange_api.CreateTypeFlangeRequest) (*flange_api.IdResponse, error) {
	id, err := h.service.CreateTypeFlange(ctx, typeFlange)
	if err != nil {
		return nil, err
	}

	return &flange_api.IdResponse{Id: id}, nil
}

func (h *FlangeHandlers) UpdateTypeFlange(ctx context.Context, typeFlange *flange_api.UpdateTypeFlangeRequest) (*flange_api.Response, error) {
	if err := h.service.UpdateTypeFlange(ctx, typeFlange); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}

func (h *FlangeHandlers) DeleteTypeFlange(ctx context.Context, typeFlange *flange_api.DeleteTypeFlangeRequest) (*flange_api.Response, error) {
	if err := h.service.DeleteTypeFlange(ctx, typeFlange); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}

func (h *FlangeHandlers) GetStandarts(ctx context.Context, req *flange_api.GetStandartsRequest) (*flange_api.StandartsResponse, error) {
	stands, err := h.service.GetStandarts(ctx, req)
	if err != nil {
		return nil, err
	}

	return &flange_api.StandartsResponse{Standarts: stands}, nil
}

func (h *FlangeHandlers) GetStandartsWithSize(ctx context.Context, req *flange_api.GetStandartsRequest,
) (*flange_api.StandartsWithSizeResponse, error) {
	stands, err := h.service.GetStandartsWithSize(ctx, req)
	if err != nil {
		return nil, err
	}

	return &flange_api.StandartsWithSizeResponse{Standarts: stands}, nil
}

func (h *FlangeHandlers) CreateStandart(ctx context.Context, stand *flange_api.CreateStandartRequest) (*flange_api.IdResponse, error) {
	id, err := h.service.CreateStandart(ctx, stand)
	if err != nil {
		return nil, err
	}

	return &flange_api.IdResponse{Id: id}, nil
}

func (h *FlangeHandlers) UpdateStandart(ctx context.Context, stand *flange_api.UpdateStandartRequest) (*flange_api.Response, error) {
	if err := h.service.UpdateStandart(ctx, stand); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}

func (h *FlangeHandlers) DeleteStandart(ctx context.Context, stand *flange_api.DeleteStandartRequest) (*flange_api.Response, error) {
	if err := h.service.DeleteStandart(ctx, stand); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}
