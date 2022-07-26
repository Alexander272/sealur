package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/user_api"
)

func (h *Handler) AddIp(ctx context.Context, ip *user_api.AddIpRequest) (*user_api.SuccessResponse, error) {
	err := h.service.IP.Add(ctx, ip)
	if err != nil {
		return nil, err
	}

	return &user_api.SuccessResponse{Success: true}, nil
}
