package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) SendInterview(ctx context.Context, req *proto.SendInterviewRequest) (*proto.SuccessResponse, error) {
	// TODO надо подключаться к фаловому и почтовому сервисам
	if err := h.service.Interview.SendInterview(req); err != nil {
		return nil, err
	}
	return &proto.SuccessResponse{Success: true}, nil
}
