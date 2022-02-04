package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetSNP(ctx context.Context, dto *proto.GetSNPRequest) (*proto.SNPResponse, error) {
	snp, err := h.service.SNP.Get(dto)
	if err != nil {
		return nil, err
	}

	return &proto.SNPResponse{Snp: snp}, nil
}

func (h *Handler) CreateSNP(ctx context.Context, dto *proto.CreateSNPRequest) (*proto.IdResponse, error) {
	snp, err := h.service.SNP.Create(dto)
	if err != nil {
		return nil, err
	}

	return snp, nil
}

func (h *Handler) UpdateSNP(ctx context.Context, dto *proto.UpdateSNPRequest) (*proto.IdResponse, error) {
	if err := h.service.SNP.Update(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}

func (h *Handler) DeleteSNP(ctx context.Context, dto *proto.DeleteSNPRequest) (*proto.IdResponse, error) {
	if err := h.service.SNP.Delete(dto); err != nil {
		return nil, err
	}
	return &proto.IdResponse{Id: dto.Id}, nil
}
