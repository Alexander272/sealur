package service

import (
	"time"

	"github.com/Alexander272/sealur/api_service/internal/repository"
	"github.com/Alexander272/sealur/api_service/pkg/auth"
)

const CookieName = "session"

type Session interface{}

type Services struct {
	Session
}

type Deps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Domain          string
}

func NewServices(deps Deps) *Services {
	return &Services{}
}
