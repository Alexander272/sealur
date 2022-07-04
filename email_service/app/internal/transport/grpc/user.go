package grpc

import (
	"context"
	"fmt"

	proto_email "github.com/Alexander272/sealur/email_service/internal/transport/grpc/proto"
)

func (h *Handler) SendConfirm(ctx context.Context, req *proto_email.ConfirmUserRequest) (*proto_email.SuccessResponse, error) {
	if err := h.service.User.SendConfirm(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to send confirm user. error: %w", err)
	}
	return &proto_email.SuccessResponse{Success: true}, nil
}

func (h *Handler) SendReject(ctx context.Context, req *proto_email.RejectUserRequest) (*proto_email.SuccessResponse, error) {
	if err := h.service.User.SendReject(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to send reject user. error: %w", err)
	}
	return &proto_email.SuccessResponse{Success: true}, nil
}

func (h *Handler) SendJoin(ctx context.Context, user *proto_email.JoinUserRequest) (*proto_email.SuccessResponse, error) {
	if err := h.service.User.SendJoin(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to send join user. error: %w", err)
	}
	return &proto_email.SuccessResponse{Success: true}, nil
}
