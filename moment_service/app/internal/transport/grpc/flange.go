package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type FlangeHandlers struct {
	service service.Flange
}

func NewFlangeHandlers(service service.Flange) *FlangeHandlers {
	return &FlangeHandlers{
		service: service,
	}
}

func (h *FlangeHandlers) CalculateFlange(ctx context.Context, req *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error) {
	// default TipF1 = 0
	res, err := h.service.Calculation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
