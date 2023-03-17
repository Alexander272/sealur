package repository

import (
	"context"

	"github.com/Alexander272/sealur/api_service/internal/config"
	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/go-redis/redis/v8"
)

type Session interface {
	Get(ctx context.Context, sessionName string) (SessionData, error)
	GetDel(ctx context.Context, sessionName string) (SessionData, error)
	Create(ctx context.Context, sessionName string, data SessionData) error
	Remove(ctx context.Context, sessionName string) error
}

type Limit interface {
	Get(ctx context.Context, clientIP string) (models.LimitData, error)
	Create(ctx context.Context, clientIP string) error
	AddAttempt(ctx context.Context, clientIP string) error
	Remove(ctx context.Context, clientIP string) error
}

type Confirm interface {
	Get(ctx context.Context, code string) (data models.ConfirmData, err error)
	Create(ctx context.Context, data models.ConfirmData) error
}

type Repositories struct {
	Session
	Limit
	Confirm
}

func NewRepo(client redis.Cmdable, conf config.AuthConfig) *Repositories {
	return &Repositories{
		Session: NewSessionRepo(client),
		Limit:   NewLimitRepo(client, conf.LimitAuthTTL),
		Confirm: NewConfirmRepo(client),
	}
}
