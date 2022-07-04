package service

import (
	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/models"
	proto_email "github.com/Alexander272/sealur/email_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/email_service/pkg/email"
)

type TestwService struct {
	sender email.Sender
	conf   config.RecipientsConfig
}

func NewTestService(sender email.Sender, conf config.RecipientsConfig) *TestwService {
	return &TestwService{
		sender: sender,
		conf:   conf,
	}
}

func (s *TestwService) SendEmail(req *proto_email.SendTestRequest) error {
	input := email.SendEmailInput{
		Subject: "Testing Join Email Sealur Pro",
		To:      []string{s.conf.Test},
	}

	switch req.Type.String() {
	case "join":
		input.Subject = "Testing Join Email Sealur Pro"

		data := models.JoinTemplate{
			Name:     "Alex",
			Login:    "testname",
			Password: "qwerty",
			Services: "сервису Sealur Pro",
			Link:     s.conf.Link,
			Email:    s.conf.Support,
		}

		if err := input.GenerateBodyFromHTML("join.html", data); err != nil {
			return err
		}

	case "confirm":
		input.Subject = "Testing Confirm Email Sealur Pro"

		data := models.ConfirmTemplate{
			Name:         "Alex",
			Organization: "Sealur",
			Position:     "developer",
			Link:         s.conf.Link,
		}

		if err := input.GenerateBodyFromHTML("confirm.html", data); err != nil {
			return err
		}

	case "interview":
		input.Subject = "Testing Interview Email Sealur Pro"

		data := proto_email.User{
			Name:         "Alex",
			Organization: "Sealur",
			Position:     "developer",
			Email:        "test@mail.com",
			City:         "Perm",
			Phone:        "-",
		}

		if err := input.GenerateBodyFromHTML("interview.html", data); err != nil {
			return err
		}

	case "order":
		input.Subject = "Testing Order Email Sealur Pro"

		data := proto_email.User{
			Name:         "Alex",
			Organization: "Sealur",
			Position:     "developer",
			Email:        "test@mail.com",
			City:         "Perm",
			Phone:        "-",
		}

		if err := input.GenerateBodyFromHTML("order.html", data); err != nil {
			return err
		}

	case "reject":
		input.Subject = "Testing Reject Email Sealur Pro"

		data := models.RejectTemplate{
			Name:  "Alex",
			Email: s.conf.Support,
		}

		if err := input.GenerateBodyFromHTML("reject.html", data); err != nil {
			return err
		}

	default:
		input.Subject = "Testing Join Email Sealur Pro"

		data := models.JoinTemplate{
			Name:     "Alex",
			Login:    "testname",
			Password: "qwerty",
			Services: "сервису Sealur Pro",
			Link:     s.conf.Link,
			Email:    s.conf.Support,
		}

		if err := input.GenerateBodyFromHTML("join.html", data); err != nil {
			return err
		}
	}

	return s.sender.Send(input)
}
