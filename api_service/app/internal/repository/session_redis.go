package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Alexander272/sealur/api_service/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type SessionData struct {
	UserId string
	Role   string
	Ua     string
	Ip     string
	Exp    time.Duration
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

func (r *SessionRepo) CreateSession(ctx context.Context, token string, data SessionData) error {
	if err := r.client.Set(ctx, token, data, data.Exp).Err(); err != nil {
		return err
	}
	return nil
}

func (r *SessionRepo) GetDelSession(ctx context.Context, token string) (data SessionData, err error) {
	cmd := r.client.GetDel(ctx, token)
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

func (r *SessionRepo) RemoveSession(ctx context.Context, token string) error {
	if err := r.client.Del(ctx, token).Err(); err != nil {
		return err
	}
	return nil
}
