package grpc

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/email_api"
)

func (h *Handler) ConfirmUser(ctx context.Context, req *email_api.ConfirmUserRequest) (*email_api.SuccessResponse, error) {
	if err := h.service.User.SendConfirmNew(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to send confirm user. error: %w", err)
	}
	return &email_api.SuccessResponse{Success: true}, nil
}

func (h *Handler) SendConfirm(ctx context.Context, req *email_api.ConfirmUserRequestOld) (*email_api.SuccessResponse, error) {
	if err := h.service.User.SendConfirm(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to send confirm user. error: %w", err)
	}
	return &email_api.SuccessResponse{Success: true}, nil
}

func (h *Handler) SendReject(ctx context.Context, req *email_api.RejectUserRequest) (*email_api.SuccessResponse, error) {
	if err := h.service.User.SendReject(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to send reject user. error: %w", err)
	}
	return &email_api.SuccessResponse{Success: true}, nil
}

func (h *Handler) SendJoin(ctx context.Context, user *email_api.JoinUserRequest) (*email_api.SuccessResponse, error) {
	if err := h.service.User.SendJoin(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to send join user. error: %w", err)
	}
	return &email_api.SuccessResponse{Success: true}, nil
}
