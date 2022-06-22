package middleware

import (
	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/service"
)

type Middleware struct {
	CookieName string
	services   *service.Services
	auth       config.AuthConfig
}

func NewMiddleware(services *service.Services, auth config.AuthConfig) *Middleware {
	return &Middleware{
		services: services,
		auth:     auth,
	}
}
