package repository

import (
	"context"
)

type Session interface {
	CreateSession(ctx context.Context, token string, data SessionData) error
	GetDelSession(ctx context.Context, token string) (SessionData, error)
	RemoveSession(ctx context.Context, token string) error
}

type Repositories struct {
	Session
}

// func NewRepo(client redis.Cmdable) *Repositories {
func NewRepo() *Repositories {
	return &Repositories{
		// Session: NewSessionRepo(client),
	}
}
