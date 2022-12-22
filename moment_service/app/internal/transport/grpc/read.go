package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment/read_api"
)

type ReadHandlers struct {
	service service.Read
	read_api.UnimplementedReadServiceServer
}

func NewReadHandlers(service service.Read) *ReadHandlers {
	return &ReadHandlers{
		service: service,
	}
}

func (h *ReadHandlers) GetFlange(ctx context.Context, req *read_api.GetFlangeRequest) (*read_api.GetFlangeResponse, error) {
	res, err := h.service.GetFlange(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *ReadHandlers) GetFloat(ctx context.Context, req *read_api.GetFloatRequest) (*read_api.GetFloatResponse, error) {
	res, err := h.service.GetFloat(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *ReadHandlers) GetDevCooling(ctx context.Context, req *read_api.GetDevCoolingtRequest) (*read_api.GetDevCoolingResponse, error) {
	res, err := h.service.GetDevCooling(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *ReadHandlers) GetAVO(ctx context.Context, req *read_api.GetAVORequest) (*read_api.GetAVOResponse, error) {
	res, err := h.service.GetAVO(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
