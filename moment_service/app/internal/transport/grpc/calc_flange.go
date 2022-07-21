package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

type CalcFlangeHandlers struct {
	service service.CalcFlange
}

func NewCalcFlangeHandlers(service service.CalcFlange) *CalcFlangeHandlers {
	return &CalcFlangeHandlers{
		service: service,
	}
}

func (h *CalcFlangeHandlers) CalculateFlange(ctx context.Context, req *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error) {
	res, err := h.service.Calculation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
