package repository

import (
	"context"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/go-redis/redis/v8"
)

type Session interface {
	Create(ctx context.Context, sessionName string, data SessionData) error
	Get(ctx context.Context, sessionName string) (SessionData, error)
	GetDel(ctx context.Context, sessionName string) (SessionData, error)
	Remove(ctx context.Context, sessionName string) error
}

type Limit interface {
	Create(ctx context.Context, clientIP string) error
	Get(ctx context.Context, clientIP string) (models.LimitData, error)
	AddAttempt(ctx context.Context, clientIP string) error
	Remove(ctx context.Context, clientIP string) error
}

type Repositories struct {
	Session
	Limit
}

func NewRepo(client redis.Cmdable, conf config.AuthConfig) *Repositories {
	return &Repositories{
		Session: NewSessionRepo(client),
		Limit:   NewLimitRepo(client, conf.LimitAuthTTL),
	}
}
