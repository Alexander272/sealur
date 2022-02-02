package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetStands(ctx context.Context, proto proto.GetStands) (stands []proto.Stand, err error) {
	stands, err = h.service.GetAll(proto)
	if err != nil {
		return nil, err
	}

	return stands, nil
}

func (h *Handler) CreateStand(ctx context.Context, proto proto.CreateStand) {

}
