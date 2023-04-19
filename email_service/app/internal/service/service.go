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
	SendNotification(context.Context, *email_api.NotificationData) error
	SendRedirect(context.Context, *email_api.RedirectData) error
}

type User interface {
	SendConfirm(context.Context, *email_api.ConfirmUserRequestOld) error
	SendConfirmNew(context.Context, *email_api.ConfirmUserRequest) error
	SendReject(context.Context, *email_api.RejectUserRequest) error
	SendJoin(context.Context, *email_api.JoinUserRequest) error
	SendRecovery(context.Context, *email_api.RecoveryPassword) error
}

type Blocked interface {
	SendBlocked(context.Context, *email_api.BlockedUserRequest) error
}

type Connect interface {
	SendFeedback(context.Context, *email_api.Feedback) error
}

type Test interface {
	SendEmail(*email_api.SendTestRequest) error
}

type Services struct {
	Interview
	Order
	User
	Connect
	Blocked
	Test
}

func NewServices(EmailSender email.Sender, conf config.RecipientsConfig) *Services {
	return &Services{
		Interview: NewInterviewService(EmailSender, conf),
		Order:     NewOrderService(EmailSender, conf),
		User:      NewUserService(EmailSender, conf),
		Connect:   NewConnectService(EmailSender, conf),
		Blocked:   NewBlockedService(EmailSender, conf),
		Test:      NewTestService(EmailSender, conf),
	}
}
