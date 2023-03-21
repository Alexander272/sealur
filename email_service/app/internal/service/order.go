package service

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/internal/constants"
	"github.com/Alexander272/sealur/email_service/internal/models"
	"github.com/Alexander272/sealur/email_service/pkg/email"
	"github.com/Alexander272/sealur_proto/api/email_api"
)

type OrderService struct {
	sender email.Sender
	conf   config.RecipientsConfig
}

func NewOrderService(sender email.Sender, conf config.RecipientsConfig) *OrderService {
	return &OrderService{
		sender: sender,
		conf:   conf,
	}
}

func (s *OrderService) SendOrder(data *email_api.OrderData, file *bytes.Buffer) error {
	input := email.SendEmailInput{
		Subject: s.conf.OrderSubject,
		To:      []string{s.conf.Order},
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

	if err := input.GenerateBodyFromHTML(constants.OrderTemplate, data.User); err != nil {
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

func (s *OrderService) SendNotification(ctx context.Context, req *email_api.NotificationData) error {
	input := email.SendEmailInput{
		Subject: s.conf.OrderSubject,
		To:      []string{req.Email},
	}

	data := models.OrderTemplate{
		Name:     req.User.Name,
		Position: req.User.Position,
		Company:  req.User.Company,
		Address:  req.User.Address,
		Email:    req.User.Email,
		Phone:    req.User.Phone,
		Link:     fmt.Sprintf("%s/%s?action=save", s.conf.OrderLink, req.OrderId),
	}

	if err := input.GenerateBodyFromHTML(constants.OrderNewTemplate, data); err != nil {
		return err
	}

	return s.sender.Send(input)
}

func (s *OrderService) readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file in zip. err %w", err)
	}
	defer f.Close()

	return io.ReadAll(f)
}
