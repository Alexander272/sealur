package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

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

func (s *InterviewService) SendInterview(data *proto_email.InterviewData, file *bytes.Buffer) error {
	input := email.SendEmailInput{
		Subject: "Опросный лист Sealur Pro",
		To:      []string{s.conf.Interview},
	}

	if err := input.GenerateBodyFromHTML("interview.html", data.User); err != nil {
		return err
	}

	if len(data.File.Name) > 1 {
		reader := bytes.NewReader(file.Bytes())
		zipReader, err := zip.NewReader(reader, data.File.Size)
		if err != nil {
			return fmt.Errorf("failed to read zip. err %w", err)
		}

		for i, zipItem := range zipReader.File {
			f, err := s.readZipFile(zipItem)
			if err != nil {
				return fmt.Errorf("failed to read file in zip. err %w", err)
			}

			input.Files = append(input.Files, email.Files{
				Filename: data.File.Name[i],
				Blob:     f,
			})
		}
	} else {
		input.Files = append(input.Files, email.Files{
			Filename: data.File.Name[0],
			Blob:     file.Bytes(),
		})
	}

	return s.sender.Send(input)
}

func (s *InterviewService) readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file in zip. err %w", err)
	}
	defer f.Close()

	return io.ReadAll(f)
}
