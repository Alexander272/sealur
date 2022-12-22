package grpc

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/email_api"
)

func (h *Handler) SendBlocked(ctx context.Context, req *email_api.BlockedUserRequest) (*email_api.SuccessResponse, error) {
	if err := h.service.Blocked.SendBlocked(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to send confirm user. error: %w", err)
	}
	return &email_api.SuccessResponse{Success: true}, nil
}
