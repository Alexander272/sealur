package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/mounting_api"
)

type MountingHandlers struct {
	service service.Mounting
	mounting_api.UnimplementedMountingServiceServer
}

func NewMountingHandlers(service service.Mounting) *MountingHandlers {
	return &MountingHandlers{
		service: service,
	}
}

func (h *MountingHandlers) GetAll(ctx context.Context, req *mounting_api.GetAllMountings) (*mounting_api.Mountings, error) {
	mountings, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return &mounting_api.Mountings{Mounting: mountings}, nil
}

func (h *MountingHandlers) Create(ctx context.Context, mounting *mounting_api.CreateMounting) (*response_model.Response, error) {
	if err := h.service.Create(ctx, mounting); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *MountingHandlers) CreateSeveral(ctx context.Context, mountings *mounting_api.CreateSeveralMounting) (*response_model.Response, error) {
	if err := h.service.CreateSeveral(ctx, mountings); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *MountingHandlers) Update(ctx context.Context, mounting *mounting_api.UpdateMounting) (*response_model.Response, error) {
	if err := h.service.Update(ctx, mounting); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *MountingHandlers) Delete(ctx context.Context, mounting *mounting_api.DeleteMounting) (*response_model.Response, error) {
	if err := h.service.Delete(ctx, mounting); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
