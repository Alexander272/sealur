package service

import (
	"context"
	"time"

	"github.com/Alexander272/sealur/api_service/internal/repository"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/proto_user"
	"github.com/Alexander272/sealur/api_service/pkg/auth"
)

type Session interface {
	SignIn(ctx context.Context, user *proto_user.User) (token string, err error)
	SingOut(ctx context.Context, userId string) error
	// Refresh(ctx context.Context, token) (*domain.Token, error)
	CheckSession(ctx context.Context, u *proto_user.User, token string) (isRefresh bool, err error)
	TokenParse(token string) (*proto_user.User, error)
}

type Services struct {
	Session
}

type Deps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps Deps) *Services {
	return &Services{
		Session: NewSessionService(deps.Repos.Session, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
	}
}
