package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type CalcFlangeHandlers struct {
	service service.CalcFlange
	moment_api.UnimplementedCalcFlangeServiceServer
}

func NewCalcFlangeHandlers(service service.CalcFlange) *CalcFlangeHandlers {
	return &CalcFlangeHandlers{
		service: service,
	}
}

func (h *CalcFlangeHandlers) CalculateFlange(ctx context.Context, req *moment_api.CalcFlangeRequest) (*moment_api.FlangeResponse, error) {
	res, err := h.service.Calculation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
