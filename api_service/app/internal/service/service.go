package service

import (
	"context"
	"net/http"
	"time"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/repository"
	"github.com/Alexander272/sealur/api_service/pkg/auth"
)

const CookieName = "session"

type Session interface {
	SignIn(ctx context.Context, dto models.SignInUserDTO, ua, ip string) (http.Cookie, models.SessionResponse, error)
	SingOut(ctx context.Context, token string) (http.Cookie, error)
	// Refresh(ctx context.Context, token, ua, ip string) (*domain.Token, *http.Cookie, error)
	TokenParse(token string) (userId string, role string, err error)
}

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
	return &Services{
		Session: NewSessionService(deps.Repos.Session, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.Domain),
	}
}
