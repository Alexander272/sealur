package smtp

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/Alexander272/sealur/email_service/pkg/email"
	"github.com/Alexander272/sealur/email_service/pkg/logger"
)

type SMTPSender struct {
	from string
	pass string
	host string
	port string
}

func NewSMTPSender(from, pass, host, port string) (*SMTPSender, error) {
	if !email.IsEmailValid(from) {
		return nil, errors.New("invalid from email")
	}

	return &SMTPSender{from: from, pass: pass, host: host, port: port}, nil
}

func (s *SMTPSender) Send(input email.SendEmailInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	delimeter := "f46d043c813270fc6b04c2d223da"

	// message := bytes.NewBuffer(nil)
	//basic email headers
	message := fmt.Sprintf("From: %s\r\n", s.from)
	message += fmt.Sprintf("To: %s\r\n", input.To)
	message += fmt.Sprintf("Subject: %s\r\n", input.Subject)
	message += "MIME-Version: 1.0\r\n"

	if len(input.Files) != 0 {
		message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", delimeter)
	}

	//place HTML message
	message += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	message += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	message += "Content-Transfer-Encoding: 7bit\r\n"
	message += fmt.Sprintf("\r\n%s", input.Body)

	//place file
	for _, f := range input.Files {
		message += fmt.Sprintf("\r\n--%s\r\n", delimeter)
		message += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
		message += "Content-Transfer-Encoding: base64\r\n"
		message += "Content-Disposition: attachment;filename=\"" + f.Filename + "\"\r\n"
		message += "\r\n" + base64.StdEncoding.EncodeToString(f.Blob)
	}

	// Authentication.
	auth := smtp.PlainAuth("", s.from, s.pass, s.host)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", s.host, s.port), auth, s.from, input.To, []byte(message))
	if err != nil {
		logger.Errorf("failed to send email. err: %w", err)
		return fmt.Errorf("failed to send email. err: %w", err)
	}

	// msg := gomail.NewMessage()
	// msg.SetHeader("From", s.from)
	// msg.SetHeader("To", input.To)
	// msg.SetHeader("Subject", input.Subject)
	// msg.SetBody("text/html", input.Body)
	// if len(input.File) != 0 {
	// 	msg.Attach(string(input.File))
	// }

	// dialer := gomail.NewDialer(s.host, s.port, s.from, s.pass)
	// if err := dialer.DialAndSend(msg); err != nil {
	// 	return fmt.Errorf("failed to sent email via smtp. err %w", err)
	// }

	return nil
}
