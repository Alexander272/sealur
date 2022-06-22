package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/transport/http/v1/proto/proto_user"
	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type SessionData struct {
	UserId       string
	AccessToken  string
	RefreshToken string
	Roles        []*proto_user.Role
	Exp          time.Duration
}

func (i SessionData) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func UnMarshalBinary(str string) SessionData {
	var data SessionData
	json.Unmarshal([]byte(str), &data)
	return data
}

type SessionRepo struct {
	client redis.Cmdable
}

func NewSessionRepo(client redis.Cmdable) *SessionRepo {
	return &SessionRepo{
		client: client,
	}
}

func (r *SessionRepo) Create(ctx context.Context, sessionName string, data SessionData) error {
	if err := r.client.Set(ctx, sessionName, data, data.Exp).Err(); err != nil {
		return err
	}
	return nil
}

func (r *SessionRepo) Get(ctx context.Context, sessionName string) (data SessionData, err error) {
	cmd := r.client.Get(ctx, sessionName)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return data, models.ErrSessionEmpty
		}
		logger.Error(cmd.Err())
		return data, cmd.Err()
	}

	str, err := cmd.Result()
	if err != nil {
		return data, err
	}
	return UnMarshalBinary(str), nil
}

func (r *SessionRepo) GetDel(ctx context.Context, sessionName string) (data SessionData, err error) {
	cmd := r.client.GetDel(ctx, sessionName)
	if cmd.Err() != nil {
		logger.Debug(cmd.Err())
		return data, cmd.Err()
	}

	str, err := cmd.Result()
	if err != nil {
		return data, err
	}
	return UnMarshalBinary(str), nil
}

func (r *SessionRepo) Remove(ctx context.Context, sessionName string) error {
	if err := r.client.Del(ctx, sessionName).Err(); err != nil {
		return err
	}
	return nil
}
