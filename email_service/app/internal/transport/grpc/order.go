package grpc

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Alexander272/sealur/email_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
)

func (h *Handler) SendOrder(stream email_api.EmailService_SendOrderServer) error {
	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("cannot receive image info %w", err)
	}

	data := req.GetData()
	file := bytes.Buffer{}

	for {
		logger.Debug("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			logger.Debug("no more data")
			break
		}

		if err != nil {
			return fmt.Errorf("cannot receive chunk data: %w", err)
		}

		chunk := req.GetFile().Content

		_, err = file.Write(chunk)
		if err != nil {
			return fmt.Errorf("cannot write chunk data: %w", err)
		}
	}

	if err := h.service.Order.SendOrder(data, &file); err != nil {
		return fmt.Errorf("failed to send order email. err: %w", err)
	}

	return stream.SendAndClose(&email_api.SuccessResponse{Success: true})
}

func (h *Handler) SendNotification(ctx context.Context, data *email_api.NotificationData) (*email_api.SuccessResponse, error) {
	if err := h.service.Order.SendNotification(ctx, data); err != nil {
		return nil, fmt.Errorf("failed to send notification to manager. error: %w", err)
	}
	return &email_api.SuccessResponse{Success: true}, nil
}

func (h *Handler) SendRedirect(ctx context.Context, data *email_api.RedirectData) (*email_api.SuccessResponse, error) {
	if err := h.service.Order.SendRedirect(ctx, data); err != nil {
		return nil, fmt.Errorf("failed to send redirect order. error: %w", err)
	}
	return &email_api.SuccessResponse{Success: true}, nil
}
