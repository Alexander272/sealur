package grpc

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type CalcCapHandlers struct {
	service service.CalcCap
	moment_api.UnimplementedCalcCapServiceServer
}

func NewCalcCapHandlers(service service.CalcCap) *CalcCapHandlers {
	return &CalcCapHandlers{
		service: service,
	}
}

func (h *CalcCapHandlers) CalculateCap(ctx context.Context, req *moment_api.CalcCapRequest) (*moment_api.CapResponse, error) {
	res, err := h.service.Calculation(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
