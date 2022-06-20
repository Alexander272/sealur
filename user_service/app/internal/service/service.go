package service

import (
	"github.com/Alexander272/sealur/user_service/internal/repo"
	proto_email "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto/email"
	"github.com/Alexander272/sealur/user_service/pkg/hasher"
)

type Services struct {
}

type Deps struct {
	Repos  *repo.Repo
	Email  proto_email.EmailServiceClient
	Hasher hasher.PasswordHasher
}

func NewServices(deps Deps) *Services {
	return &Services{}
}
