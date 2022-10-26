package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
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

func (h *CalcHandlers) CalculateFloat(ctx context.Context, req *calc_api.FloatRequest) (*calc_api.FloatResponse, error) {
	res, err := h.service.CalculationFloat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *CalcHandlers) CalculateDevCooling(ctx context.Context, req *calc_api.DevCoolingRequest) (*calc_api.DevCoolingResponse, error) {
	res, err := h.service.CalculateDevCooling(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *CalcHandlers) CalculateExCircle(ctx context.Context, req *calc_api.ExpressCircleRequest) (*calc_api.ExpressCircleResponse, error) {
	res, err := h.service.CalculateExCircle(ctx, req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return res, nil
}
