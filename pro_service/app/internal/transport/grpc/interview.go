package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) SendInterview(ctx context.Context, req *proto.SendInterviewRequest) (*proto.SuccessResponse, error) {
	if err := h.service.Interview.SendInterview(ctx, req); err != nil {
		return nil, err
	}
	return &proto.SuccessResponse{Success: true}, nil
}
