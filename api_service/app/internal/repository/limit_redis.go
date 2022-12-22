package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/go-redis/redis/v8"
)

func (r *LimitRepo) UnMarshalBinary(str string) models.LimitData {
	var data models.LimitData
	json.Unmarshal([]byte(str), &data)
	return data
}

type LimitRepo struct {
	client   redis.Cmdable
	limitTTL time.Duration
}

func NewLimitRepo(client redis.Cmdable, limitTTL time.Duration) *LimitRepo {
	return &LimitRepo{
		client:   client,
		limitTTL: limitTTL,
	}
}

func (r *LimitRepo) Create(ctx context.Context, clientIP string) error {
	data := models.LimitData{
		ClientIP: clientIP,
		Count:    1,
		Exp:      r.limitTTL,
	}

	if err := r.client.Set(ctx, clientIP, data, data.Exp).Err(); err != nil {
		return err
	}
	return nil
}

func (r *LimitRepo) Get(ctx context.Context, clientIP string) (data models.LimitData, err error) {
	cmd := r.client.Get(ctx, clientIP)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return data, models.ErrClientIPNotFound
		}
		return data, cmd.Err()
	}

	str, err := cmd.Result()
	if err != nil {
		return data, err
	}
	return r.UnMarshalBinary(str), nil
}

func (r *LimitRepo) AddAttempt(ctx context.Context, clientIP string) error {
	cmd := r.client.Get(ctx, clientIP)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return models.ErrClientIPNotFound
		}
		return cmd.Err()
	}
	str, err := cmd.Result()
	if err != nil {
		return err
	}

	data := r.UnMarshalBinary(str)
	data.Count += 1
	data.Exp = r.limitTTL
	if err := r.client.Set(ctx, clientIP, data, data.Exp).Err(); err != nil {
		return err
	}
	return nil
}

func (r *LimitRepo) Remove(ctx context.Context, clientIP string) error {
	if err := r.client.Del(ctx, clientIP).Err(); err != nil {
		return err
	}
	return nil
}
