package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetSNP(ctx context.Context, dto *pro_api.GetSNPRequest) (*pro_api.SNPResponse, error) {
	snp, err := h.service.SNP.Get(dto)
	if err != nil {
		return nil, err
	}

	return &pro_api.SNPResponse{Snp: snp}, nil
}

func (h *Handler) CreateSNP(ctx context.Context, dto *pro_api.CreateSNPRequest) (*pro_api.IdResponse, error) {
	snp, err := h.service.SNP.Create(dto)
	if err != nil {
		return nil, err
	}

	return snp, nil
}

func (h *Handler) UpdateSNP(ctx context.Context, dto *pro_api.UpdateSNPRequest) (*pro_api.IdResponse, error) {
	if err := h.service.SNP.Update(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteSNP(ctx context.Context, dto *pro_api.DeleteSNPRequest) (*pro_api.IdResponse, error) {
	if err := h.service.SNP.Delete(dto); err != nil {
		return nil, err
	}
	return &pro_api.IdResponse{Id: dto.Id}, nil
}
