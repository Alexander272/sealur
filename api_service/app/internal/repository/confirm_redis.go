package repository

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/go-redis/redis/v8"
)

type ConfirmRepo struct {
	client redis.Cmdable
}

func NewConfirmRepo(client redis.Cmdable) *ConfirmRepo {
	return &ConfirmRepo{
		client: client,
	}
}

func (r *ConfirmRepo) Get(ctx context.Context, code string) (data models.ConfirmData, err error) {
	cmd := r.client.Get(ctx, code)
	if cmd.Err() != nil {
		return models.ConfirmData{}, fmt.Errorf("failed to execute query. error: %w", cmd.Err())
	}

	str, err := cmd.Result()
	if err != nil {
		return models.ConfirmData{}, fmt.Errorf("failed to get result. error: %w", err)
	}

	data.UnMarshalBinary(str)

	return data, nil
}

func (r *ConfirmRepo) Create(ctx context.Context, data models.ConfirmData) error {
	if err := r.client.Set(ctx, data.Code, data, data.Exp).Err(); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
