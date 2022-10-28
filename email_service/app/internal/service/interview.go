package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/constants"
	"github.com/Alexander272/sealur/email_service/pkg/email"
	"github.com/Alexander272/sealur_proto/api/email_api"
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

func (s *InterviewService) SendInterview(data *email_api.InterviewData, file *bytes.Buffer) error {
	input := email.SendEmailInput{
		Subject: s.conf.InterviewSubject,
		To:      []string{s.conf.Interview},
	}

	if data.User.City == "" {
		data.User.City = "-"
	}
	if data.User.Position == "" {
		data.User.Position = "-"
	}
	if data.User.Phone == "" {
		data.User.Phone = "-"
	}

	if err := input.GenerateBodyFromHTML(constants.InterviewTemplate, data.User); err != nil {
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
