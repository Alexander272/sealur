package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
)

type CalcHandlers struct {
	service service.Calc
	calc_api.UnimplementedCalcServiceServer
}

func NewCalcHandlers(service service.Calc) *CalcHandlers {
	return &CalcHandlers{
		service: service,
	}
}

func (h *CalcHandlers) CalculateFlange(ctx context.Context, req *calc_api.FlangeRequest) (*calc_api.FlangeResponse, error) {
	res, err := h.service.CalculationFlange(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *CalcHandlers) CalculateCap(ctx context.Context, req *calc_api.CapRequest) (*calc_api.CapResponse, error) {
	res, err := h.service.CalculationCap(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
