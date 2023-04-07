package grpc

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/email_api"
)

func (h *Handler) SendFeedback(ctx context.Context, feedback *email_api.Feedback) (*email_api.SuccessResponse, error) {
	if err := h.service.Connect.SendFeedback(ctx, feedback); err != nil {
		return nil, fmt.Errorf("failed to send feedback. error: %w", err)
	}
	return &email_api.SuccessResponse{Success: true}, nil
}
