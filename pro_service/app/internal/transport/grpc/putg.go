package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/putg_api"
)

type PutgHandlers struct {
	service service.Putg
	putg_api.UnimplementedPutgDataServiceServer
}

func NewPutgHandlers(service service.Putg) *PutgHandlers {
	return &PutgHandlers{
		service: service,
	}
}

func (h *PutgHandlers) GetBase(ctx context.Context, req *putg_api.GetPutgBase) (*putg_api.PutgBase, error) {
	putgBase, err := h.service.GetBase(ctx, req)
	if err != nil {
		return nil, err
	}
	return putgBase, nil
}

func (h *PutgHandlers) GetData(ctx context.Context, req *putg_api.GetPutgData) (*putg_api.PutgData, error) {
	putgData, err := h.service.GetData(ctx, req)
	if err != nil {
		return nil, err
	}
	return putgData, nil
}

func (h *PutgHandlers) Get(ctx context.Context, req *putg_api.GetPutg) (*putg_api.Putg, error) {
	putg, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return putg, nil
}

// func (h *Handler) GetPutg(ctx context.Context, dto *pro_api.GetPutgRequest) (*pro_api.PutgResponse, error) {
// 	putg, err := h.service.Putg.Get(dto)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &pro_api.PutgResponse{Putg: putg}, nil
// }

// func (h *Handler) CreatePutg(ctx context.Context, dto *pro_api.CreatePutgRequest) (*pro_api.IdResponse, error) {
// 	putg, err := h.service.Putg.Create(dto)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return putg, nil
// }

// func (h *Handler) UpdatePutg(ctx context.Context, dto *pro_api.UpdatePutgRequest) (*pro_api.IdResponse, error) {
// 	if err := h.service.Putg.Update(dto); err != nil {
// 		return nil, err
// 	}
// 	return &pro_api.IdResponse{Id: dto.Id}, nil
// }

// func (h *Handler) DeletePutg(ctx context.Context, dto *pro_api.DeletePutgRequest) (*pro_api.IdResponse, error) {
// 	if err := h.service.Putg.Delete(dto); err != nil {
// 		return nil, err
// 	}
// 	return &pro_api.IdResponse{Id: dto.Id}, nil
// }
