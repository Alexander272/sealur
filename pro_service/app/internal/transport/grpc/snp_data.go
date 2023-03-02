package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_data_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_data_api"
)

type SnpDataHandlers struct {
	service service.SnpData
	snp_data_api.UnimplementedSnpDataServiceServer
}

func NewSnpDataHandlers(service service.SnpData) *SnpDataHandlers {
	return &SnpDataHandlers{
		service: service,
	}
}

func (h *SnpDataHandlers) Get(ctx context.Context, req *snp_data_api.GetSnpData) (*snp_data_api.SnpData, error) {
	snp, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	// TODO переделать proto, должен возвращаться объект, а не массив
	return &snp_data_api.SnpData{Snp: []*snp_data_model.SnpData{snp}}, nil
}

func (h *SnpDataHandlers) Create(ctx context.Context, snp *snp_data_api.CreateSnpData) (*response_model.Response, error) {
	if err := h.service.Create(ctx, snp); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpDataHandlers) Update(ctx context.Context, snp *snp_data_api.UpdateSnpData) (*response_model.Response, error) {
	if err := h.service.Update(ctx, snp); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *SnpDataHandlers) Delete(ctx context.Context, snp *snp_data_api.DeleteSnpData) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, snp); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
