package service

import (
	"bytes"
	"context"

	"github.com/Alexander272/sealur/email_service/internal/config"
	"github.com/Alexander272/sealur/email_service/pkg/email"
	"github.com/Alexander272/sealur_proto/api/email_api"
)

type Interview interface {
	SendInterview(*email_api.InterviewData, *bytes.Buffer) error
}

type Order interface {
	SendOrder(*email_api.OrderData, *bytes.Buffer) error
}

type User interface {
	SendConfirm(ctx context.Context, req *email_api.ConfirmUserRequest) error
	SendReject(ctx context.Context, req *email_api.RejectUserRequest) error
	SendJoin(ctx context.Context, user *email_api.JoinUserRequest) error
}

type Blocked interface {
	SendBlocked(ctx context.Context, req *email_api.BlockedUserRequest) error
}

type Test interface {
	SendEmail(*email_api.SendTestRequest) error
}

type Services struct {
	Interview
	Order
	User
	Test
	Blocked
}

func NewServices(EmailSender email.Sender, conf config.RecipientsConfig) *Services {
	return &Services{
		Interview: NewInterviewService(EmailSender, conf),
		Order:     NewOrderService(EmailSender, conf),
		User:      NewUserService(EmailSender, conf),
		Blocked:   NewBlockedService(EmailSender, conf),
		Test:      NewTestService(EmailSender, conf),
	}
}
