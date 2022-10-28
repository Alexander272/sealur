package service

import (
	"context"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/constants"
	"github.com/Alexander272/sealur/email_service/internal/models"
	"github.com/Alexander272/sealur/email_service/pkg/email"
	"github.com/Alexander272/sealur_proto/api/email_api"
)

type BlockedService struct {
	sender email.Sender
	conf   config.RecipientsConfig
}

func NewBlockedService(sender email.Sender, conf config.RecipientsConfig) *BlockedService {
	return &BlockedService{
		sender: sender,
		conf:   conf,
	}
}

func (s *BlockedService) SendBlocked(ctx context.Context, req *email_api.BlockedUserRequest) error {
	input := email.SendEmailInput{
		Subject: s.conf.BlockedSubject,
		To:      []string{s.conf.Blocked},
	}

	data := models.BlockedTemplate{
		Ip:    req.Ip,
		Login: req.Login,
	}

	if err := input.GenerateBodyFromHTML(constants.BlockedTemplate, data); err != nil {
		return err
	}

	return s.sender.Send(input)
}
