package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/repository"
	"github.com/Alexander272/sealur/api_service/pkg/auth"
)

type SessionService struct {
	repo            repository.Session
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	domain          string
}

func NewSessionService(repo repository.Session, manager auth.TokenManager, accessTTL, refreshTTL time.Duration, domain string) *SessionService {
	return &SessionService{
		repo:            repo,
		tokenManager:    manager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
		domain:          domain,
	}
}

func (s *SessionService) SignIn(ctx context.Context, dto models.SignInUserDTO, ua, ip string) (http.Cookie, models.SessionResponse, error) {
	//TODO сделать запрос на сервис пользователей
	// TODO изменить secure на true

	// iat, accessToken, err := s.tokenManager.NewJWT(user.Id, user.Email, user.Role, s.accessTokenTTL)
	// if err != nil {
	// 	return http.Cookie{}, models.SessionResponse{}, err
	// }
	// refreshToken, err := s.tokenManager.NewRefreshToken()
	// if err != nil {
	// 	return http.Cookie{}, models.SessionResponse{}, err
	// }

	// data := repository.SessionData{
	// 	UserId: user.Id,
	// 	Role:   user.Role,
	// 	Ua:     ua,
	// 	Ip:     ip,
	// 	Exp:    s.refreshTokenTTL,
	// }

	// if err := s.repo.CreateSession(ctx, refreshToken, data); err != nil {
	// 	return nil, nil, err
	// }

	// cookie := http.Cookie{
	// 	Name:     CookieName,
	// 	Value:    refreshToken,
	// 	MaxAge:   int(s.refreshTokenTTL.Seconds()),
	// 	Path:     "/",
	// 	Domain:   s.domain,
	// 	Secure:   false,
	// 	HttpOnly: true,
	// }

	// token := models.SessionResponse{
	// 	Token: models.Token{
	// 		AccessToken: accessToken,
	// 		Exp:         iat.Add(s.accessTokenTTL).Unix(),
	// 	},
	// 	UserId: user.Id,
	// 	Role:   user.Role,
	// }

	// return cookie, token, nil
	return http.Cookie{}, models.SessionResponse{}, errors.New("not implemented")
}

func (s *SessionService) SingOut(ctx context.Context, token string) (http.Cookie, error) {
	cookie := http.Cookie{
		Name:     CookieName,
		Value:    "",
		MaxAge:   1,
		Path:     "/",
		Domain:   s.domain,
		Secure:   false,
		HttpOnly: true,
	}

	err := s.repo.RemoveSession(ctx, token)
	if err != nil {
		return cookie, err
	}

	return cookie, nil
}

//TODO дописать refresh
// func (s *SessionService) Refresh(ctx context.Context, token, ua, ip string) (*domain.Token, *http.Cookie, error)

func (s *SessionService) TokenParse(token string) (userId string, role string, err error) {
	claims, err := s.tokenManager.Parse(token)
	if err != nil {
		return "", "", err
	}
	return claims["userId"].(string), claims["role"].(string), err
}
