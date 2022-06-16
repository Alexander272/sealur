package email

import (
	"bytes"
	"errors"
	"html/template"
	"path"

	"github.com/Alexander272/sealur/email_service/pkg/logger"
)

type SendEmailInput struct {
	To      []string
	Subject string
	Body    string
	Files   []Files
}

type Files struct {
	Filename string
	Blob     []byte
}

type Sender interface {
	Send(input SendEmailInput) error
}

func (e *SendEmailInput) GenerateBodyFromHTML(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(path.Join("template", templateFileName))
	if err != nil {
		logger.Errorf("failed to parse file %s:%s", templateFileName, err.Error())

		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	e.Body = buf.String()

	return nil
}

func (e *SendEmailInput) Validate() error {
	if len(e.To) == 0 {
		return errors.New("empty to")
	}

	if e.Subject == "" || e.Body == "" {
		return errors.New("empty subject/body")
	}

	for _, v := range e.To {
		if !IsEmailValid(v) {
			return errors.New("invalid to email")
		}
	}

	return nil
}
