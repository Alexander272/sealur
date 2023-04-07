package service

import (
	"context"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/constants"
	"github.com/Alexander272/sealur/email_service/pkg/email"
	"github.com/Alexander272/sealur_proto/api/email_api"
)

type ConnectService struct {
	sender email.Sender
	conf   config.RecipientsConfig
}

func NewConnectService(sender email.Sender, conf config.RecipientsConfig) *ConnectService {
	return &ConnectService{
		sender: sender,
		conf:   conf,
	}
}

func (s *ConnectService) SendFeedback(ctx context.Context, feedback *email_api.Feedback) error {
	input := email.SendEmailInput{
		Subject: feedback.Subject,
		To:      []string{s.conf.Connect},
	}

	if err := input.GenerateBodyFromHTML(constants.FeedbackTemplate, feedback); err != nil {
		return err
	}

	return s.sender.Send(input)
}
