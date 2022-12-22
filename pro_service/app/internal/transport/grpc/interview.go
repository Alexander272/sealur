package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) SendInterview(ctx context.Context, req *pro_api.SendInterviewRequest) (*pro_api.SuccessResponse, error) {
	if err := h.service.Interview.SendInterview(ctx, req); err != nil {
		return nil, err
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}
