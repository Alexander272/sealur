package service

import (
	"bytes"

	"github.com/Alexander272/sealur/email_service/internal/config"
	proto_email "github.com/Alexander272/sealur/email_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/email_service/pkg/email"
)

type InterviewService struct {
	sender email.Sender
	conf   config.RecipientsConfig
}

func NewInterviewService(sender email.Sender, conf config.RecipientsConfig) *InterviewService {
	return &InterviewService{
		sender: sender,
		conf:   conf,
	}
}

func (s *InterviewService) SendInterview(data *proto_email.InterviewData, file bytes.Buffer) error {
	input := email.SendEmailInput{
		Subject: "Опросный лист",
		To:      []string{s.conf.Interview},
	}

	if err := input.GenerateBodyFromHTML("interview.html", data.User); err != nil {
		return err
	}

	input.Files = append(input.Files, email.Files{
		Filename: data.File.Name,
		Blob:     file.Bytes(),
	})

	return s.sender.Send(input)
}
