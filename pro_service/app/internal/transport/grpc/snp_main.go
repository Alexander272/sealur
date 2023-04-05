package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/snp_api"
)

type SnpHandlers struct {
	service service.Snp
	snp_api.UnimplementedSnpDataServiceServer
}

func NewSnpHandlers(service service.Snp) *SnpHandlers {
	return &SnpHandlers{
		service: service,
	}
}

func (h *SnpHandlers) Get(ctx context.Context, req *snp_api.GetSnp) (*snp_api.Snp, error) {
	snp, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return snp, nil
}

func (h *SnpHandlers) GetData(ctx context.Context, req *snp_api.GetSnpData) (*snp_api.SnpData, error) {
	snpData, err := h.service.GetData(ctx, req)
	if err != nil {
		return nil, err
	}
	return &snp_api.SnpData{SnpData: snpData}, nil
}
