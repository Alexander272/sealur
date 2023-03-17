package service

import (
	"context"
	"time"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/repository"
	"github.com/Alexander272/sealur/api_service/pkg/auth"
	"github.com/Alexander272/sealur_proto/api/user/models/user_model"
)

type Session interface {
	SignIn(ctx context.Context, user *user_model.User) (token string, err error)
	SingOut(ctx context.Context, userId string) error
	CheckSession(ctx context.Context, u *user_model.User, token string) (isRefresh bool, err error)
	TokenParse(token string) (*user_model.User, error)
}

type Limit interface {
	Get(ctx context.Context, clientIP string) (models.LimitData, error)
	Create(ctx context.Context, clientIP string) error
	AddAttempt(ctx context.Context, clientIP string) error
	Remove(ctx context.Context, clientIP string) error
}

type Confirm interface {
	Get(ctx context.Context, code string) (models.ConfirmData, error)
	Create(ctx context.Context, userId string) (string, error)
}

type Services struct {
	Session
	Limit
	Confirm
}

type Deps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	ConfirmTTL      time.Duration
}

func NewServices(deps Deps) *Services {
	return &Services{
		Session: NewSessionService(deps.Repos.Session, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Limit:   NewLimitService(deps.Repos.Limit),
		Confirm: NewConfirmService(deps.Repos.Confirm, deps.TokenManager, deps.ConfirmTTL),
	}
}
