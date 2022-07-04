package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/models"
	proto_email "github.com/Alexander272/sealur/email_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/email_service/pkg/email"
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

func (s *UserService) SendConfirm(ctx context.Context, req *proto_email.ConfirmUserRequest) error {
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

	if err := input.GenerateBodyFromHTML("confirm.html", data); err != nil {
		return err
	}

	return s.sender.Send(input)
}

func (s *UserService) SendReject(ctx context.Context, req *proto_email.RejectUserRequest) error {
	input := email.SendEmailInput{
		Subject: s.conf.RejectSubject,
		To:      []string{req.Email},
	}

	data := models.RejectTemplate{
		Name:  req.Name,
		Email: s.conf.Support,
	}

	if err := input.GenerateBodyFromHTML("reject.html", data); err != nil {
		return err
	}

	return s.sender.Send(input)
}

func (s *UserService) SendJoin(ctx context.Context, user *proto_email.JoinUserRequest) error {
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

	if err := input.GenerateBodyFromHTML("join.html", data); err != nil {
		return err
	}

	return s.sender.Send(input)
}
