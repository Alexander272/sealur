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
	res, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
