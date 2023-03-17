package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/constants"
	"github.com/Alexander272/sealur/email_service/internal/models"
	"github.com/Alexander272/sealur/email_service/pkg/email"
	"github.com/Alexander272/sealur_proto/api/email_api"
)

type UserService struct {
	sender email.Sender
	conf   config.RecipientsConfig
}

func NewUserService(sender email.Sender, conf config.RecipientsConfig) *UserService {
	return &UserService{
		sender: sender,
		conf:   conf,
	}
}

func (s *UserService) SendConfirm(ctx context.Context, req *email_api.ConfirmUserRequestOld) error {
	input := email.SendEmailInput{
		Subject: s.conf.ConfirmSubject,
		To:      []string{s.conf.Confirm},
	}

	data := models.ConfirmTemplate{
		Name:         req.Name,
		Organization: req.Organization,
		Position:     req.Position,
		Link:         s.conf.Link,
	}

	if data.Position == "" {
		data.Position = "-"
	}

	if err := input.GenerateBodyFromHTML(constants.ConfirmTemplate, data); err != nil {
		return err
	}

	return s.sender.Send(input)
}

func (s *UserService) SendConfirmNew(ctx context.Context, req *email_api.ConfirmUserRequest) error {
	input := email.SendEmailInput{
		Subject: s.conf.ConfirmSubjectNew,
		To:      []string{req.Email},
	}

	data := models.ConfirmTemplateNew{
		Name:  req.Name,
		Link:  req.Link,
		Email: s.conf.Support,
	}

	if err := input.GenerateBodyFromHTML(constants.NewConfirmTemplate, data); err != nil {
		return err
	}

	return s.sender.Send(input)
}

func (s *UserService) SendReject(ctx context.Context, req *email_api.RejectUserRequest) error {
	input := email.SendEmailInput{
		Subject: s.conf.RejectSubject,
		To:      []string{req.Email},
	}

	data := models.RejectTemplate{
		Name:  req.Name,
		Email: s.conf.Support,
	}

	if err := input.GenerateBodyFromHTML(constants.RejectTemplate, data); err != nil {
		return err
	}

	return s.sender.Send(input)
}

func (s *UserService) SendJoin(ctx context.Context, user *email_api.JoinUserRequest) error {
	input := email.SendEmailInput{
		Subject: s.conf.JoinSubject,
		To:      []string{user.Email},
	}

	data := models.JoinTemplate{
		Name:     user.Name,
		Login:    user.Login,
		Password: user.Password,
		Link:     s.conf.Link,
		Email:    s.conf.Support,
	}

	if len(user.Services) > 1 {
		data.Services = fmt.Sprintf("сервисам %s", strings.Join(user.Services, ", "))
	} else {
		data.Services = fmt.Sprintf("сервису %s", user.Services[0])
	}

	if err := input.GenerateBodyFromHTML(constants.JoinTemplate, data); err != nil {
		return err
	}

	return s.sender.Send(input)
}
