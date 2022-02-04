package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetAllAdditional(ctx context.Context, dto *proto.GetAllAddRequest) (*proto.AdditionalResponse, error) {
	addit, err := h.service.Addit.GetAll()
	if err != nil {
		return nil, err
	}

	return &proto.AdditionalResponse{Additionals: addit}, nil
}

func (h *Handler) CreateAdditional(ctx context.Context, dto *proto.CreateAddRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.Create(dto)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateMat(ctx context.Context, dto *proto.UpdateAddMatRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateMat(dto)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateMod(ctx context.Context, dto *proto.UpdateAddModRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateMod(dto)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateTemp(ctx context.Context, dto *proto.UpdateAddTemRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateTemp(dto)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateMoun(ctx context.Context, dto *proto.UpdateAddMounRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateMoun(dto)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateGrap(ctx context.Context, dto *proto.UpdateAddGrapRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.UpdateGrap(dto)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateTypeFl(ctx context.Context, dto *proto.UpdateAddTypeFlRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.UpdateTypeFl(dto)
	if err != nil {
		return nil, err
	}

	return success, nil
}
