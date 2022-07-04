package service

import (
	"bytes"
	"context"

	"github.com/Alexander272/sealur/email_service/internal/config"
	proto_email "github.com/Alexander272/sealur/email_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/email_service/pkg/email"
)

type Interview interface {
	SendInterview(*proto_email.InterviewData, *bytes.Buffer) error
}

type Order interface {
	SendOrder(*proto_email.OrderData, *bytes.Buffer) error
}

type User interface {
	SendConfirm(ctx context.Context, req *proto_email.ConfirmUserRequest) error
	SendReject(ctx context.Context, req *proto_email.RejectUserRequest) error
	SendJoin(ctx context.Context, user *proto_email.JoinUserRequest) error
}

type Test interface {
	SendEmail(*proto_email.SendTestRequest) error
}

type Services struct {
	Interview
	Order
	User
	Test
}

func NewServices(EmailSender email.Sender, conf config.RecipientsConfig) *Services {
	return &Services{
		Interview: NewInterviewService(EmailSender, conf),
		Order:     NewOrderService(EmailSender, conf),
		User:      NewUserService(EmailSender, conf),
		Test:      NewTestService(EmailSender, conf),
	}
}
