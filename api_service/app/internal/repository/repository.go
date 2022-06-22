package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Session interface {
	Create(ctx context.Context, sessionName string, data SessionData) error
	Get(ctx context.Context, sessionName string) (SessionData, error)
	GetDel(ctx context.Context, sessionName string) (SessionData, error)
	Remove(ctx context.Context, sessionName string) error
}

type Repositories struct {
	Session
}

func NewRepo(client redis.Cmdable) *Repositories {
	return &Repositories{
		Session: NewSessionRepo(client),
	}
}
