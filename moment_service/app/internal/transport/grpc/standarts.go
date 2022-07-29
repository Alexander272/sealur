package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (h *FlangeHandlers) GetTypeFlange(ctx context.Context, req *moment_api.GetTypeFlangeRequest) (*moment_api.TypeFlangeResponse, error) {
	typeFlange, err := h.service.GetTypeFlange(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_api.TypeFlangeResponse{TypeFlanges: typeFlange}, nil
}

func (h *FlangeHandlers) CreateTypeFlange(ctx context.Context, typeFlange *moment_api.CreateTypeFlangeRequest) (*moment_api.IdResponse, error) {
	id, err := h.service.CreateTypeFlange(ctx, typeFlange)
	if err != nil {
		return nil, err
	}

	return &moment_api.IdResponse{Id: id}, nil
}

func (h *FlangeHandlers) UpdateTypeFlange(ctx context.Context, typeFlange *moment_api.UpdateTypeFlangeRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateTypeFlange(ctx, typeFlange); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *FlangeHandlers) DeleteTypeFlange(ctx context.Context, typeFlange *moment_api.DeleteTypeFlangeRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteTypeFlange(ctx, typeFlange); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *FlangeHandlers) GetStandarts(ctx context.Context, req *moment_api.GetStandartsRequest) (*moment_api.StandartsResponse, error) {
	stands, err := h.service.GetStandarts(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_api.StandartsResponse{Standarts: stands}, nil
}

func (h *FlangeHandlers) GetStandartsWithSize(ctx context.Context, req *moment_api.GetStandartsRequest,
) (*moment_api.StandartsWithSizeResponse, error) {
	stands, err := h.service.GetStandartsWithSize(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_api.StandartsWithSizeResponse{Standarts: stands}, nil
}

func (h *FlangeHandlers) CreateStandart(ctx context.Context, stand *moment_api.CreateStandartRequest) (*moment_api.IdResponse, error) {
	id, err := h.service.CreateStandart(ctx, stand)
	if err != nil {
		return nil, err
	}

	return &moment_api.IdResponse{Id: id}, nil
}

func (h *FlangeHandlers) UpdateStandart(ctx context.Context, stand *moment_api.UpdateStandartRequest) (*moment_api.Response, error) {
	if err := h.service.UpdateStandart(ctx, stand); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}

func (h *FlangeHandlers) DeleteStandart(ctx context.Context, stand *moment_api.DeleteStandartRequest) (*moment_api.Response, error) {
	if err := h.service.DeleteStandart(ctx, stand); err != nil {
		return nil, err
	}
	return &moment_api.Response{}, nil
}
