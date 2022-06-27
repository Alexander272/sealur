package grpc

import (
	"context"

	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
)

func (h *Handler) AddIp(ctx context.Context, ip *proto_user.AddIpRequest) (*proto_user.SuccessResponse, error) {
	err := h.service.IP.Add(ctx, ip)
	if err != nil {
		return nil, err
	}

	return &proto_user.SuccessResponse{Success: true}, nil
}
