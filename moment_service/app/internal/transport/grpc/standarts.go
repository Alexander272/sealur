package grpc

import (
	"context"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (h *FlangeHandlers) GetTypeFlange(ctx context.Context, req *moment_proto.GetTypeFlangeRequest) (*moment_proto.TypeFlangeResponse, error) {
	typeFlange, err := h.service.GetTypeFlange(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.TypeFlangeResponse{TypeFlanges: typeFlange}, nil
}

func (h *FlangeHandlers) CreateTypeFlange(ctx context.Context, typeFlange *moment_proto.CreateTypeFlangeRequest) (*moment_proto.IdResponse, error) {
	id, err := h.service.CreateTypeFlange(ctx, typeFlange)
	if err != nil {
		return nil, err
	}

	return &moment_proto.IdResponse{Id: id}, nil
}

func (h *FlangeHandlers) UpdateTypeFlange(ctx context.Context, typeFlange *moment_proto.UpdateTypeFlangeRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateTypeFlange(ctx, typeFlange); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *FlangeHandlers) DeleteTypeFlange(ctx context.Context, typeFlange *moment_proto.DeleteTypeFlangeRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteTypeFlange(ctx, typeFlange); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *FlangeHandlers) GetStandarts(ctx context.Context, req *moment_proto.GetStandartsRequest) (*moment_proto.StandartsResponse, error) {
	stands, err := h.service.GetStandarts(ctx, req)
	if err != nil {
		return nil, err
	}

	return &moment_proto.StandartsResponse{Standarts: stands}, nil
}

func (h *FlangeHandlers) CreateStandart(ctx context.Context, stand *moment_proto.CreateStandartRequest) (*moment_proto.IdResponse, error) {
	id, err := h.service.CreateStandart(ctx, stand)
	if err != nil {
		return nil, err
	}

	return &moment_proto.IdResponse{Id: id}, nil
}

func (h *FlangeHandlers) UpdateStandart(ctx context.Context, stand *moment_proto.UpdateStandartRequest) (*moment_proto.Response, error) {
	if err := h.service.UpdateStandart(ctx, stand); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}

func (h *FlangeHandlers) DeleteStandart(ctx context.Context, stand *moment_proto.DeleteStandartRequest) (*moment_proto.Response, error) {
	if err := h.service.DeleteStandart(ctx, stand); err != nil {
		return nil, err
	}
	return &moment_proto.Response{}, nil
}
