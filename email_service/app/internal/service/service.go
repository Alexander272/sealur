package service

import (
	"bytes"

	"github.com/Alexander272/sealur/email_service/internal/config"
	proto_email "github.com/Alexander272/sealur/email_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/email_service/pkg/email"
)

type Interview interface {
	SendInterview(*proto_email.InterviewData, *bytes.Buffer) error
}

type Services struct {
	Interview
}

func NewServices(EmailSender email.Sender, conf config.RecipientsConfig) *Services {
	return &Services{
		Interview: NewInterviewService(EmailSender, conf),
	}
}