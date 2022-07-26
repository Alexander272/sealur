package grpc

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Alexander272/sealur/email_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
)

func (h *Handler) SendInterview(stream email_api.EmailService_SendInterviewServer) error {
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

	if err := h.service.Interview.SendInterview(data, &file); err != nil {
		return fmt.Errorf("failed to send interview email. err: %w", err)
	}

	return stream.SendAndClose(&email_api.SuccessResponse{Success: true})
}
