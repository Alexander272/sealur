package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
)

func (h *FlangeHandlers) GetBolts(ctx context.Context, req *flange_api.GetBoltsRequest) (*flange_api.BoltsResponse, error) {
	bolts, err := h.service.GetBolts(ctx, req)
	if err != nil {
		return nil, err
	}

	return &flange_api.BoltsResponse{Bolts: bolts}, nil
}

func (h *FlangeHandlers) GetAllBolts(ctx context.Context, req *flange_api.GetBoltsRequest) (*flange_api.BoltsResponse, error) {
	bolts, err := h.service.GetAllBolts(ctx, req)
	if err != nil {
		return nil, err
	}

	return &flange_api.BoltsResponse{Bolts: bolts}, nil
}

func (h *FlangeHandlers) CreateBolt(ctx context.Context, bolt *flange_api.CreateBoltRequest) (*flange_api.Response, error) {
	if err := h.service.CreateBolt(ctx, bolt); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}

func (h *FlangeHandlers) CreateBolts(ctx context.Context, bolt *flange_api.CreateBoltsRequest) (*flange_api.Response, error) {
	if err := h.service.CreateBolts(ctx, bolt); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}

func (h *FlangeHandlers) UpdateBolt(ctx context.Context, bolt *flange_api.UpdateBoltRequest) (*flange_api.Response, error) {
	if err := h.service.UpdateBolt(ctx, bolt); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}

func (h *FlangeHandlers) DeleteBolt(ctx context.Context, bolt *flange_api.DeleteBoltRequest) (*flange_api.Response, error) {
	if err := h.service.DeleteBolt(ctx, bolt); err != nil {
		return nil, err
	}
	return &flange_api.Response{}, nil
}
