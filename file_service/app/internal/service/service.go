package service

import (
	"context"

	"github.com/Alexander272/sealur/file_service/internal/models"
	"github.com/Alexander272/sealur/file_service/internal/repository"
)

type Store interface {
	GetFile(ctx context.Context, bucket, group, id, name string) (*models.File, error)
	GetFilesByGroup(ctx context.Context, bucket, group string) ([]*models.File, error)
	Create(ctx context.Context, bucket string, dto models.CreateFileDTO) (string, error)
	Copy(ctx context.Context, bucket, group, newGroup, id string) error
	CopyGroup(ctx context.Context, bucket, group, newGroup string) error
	Delete(ctx context.Context, bucket, group, id, name string) error
	DeleteGroup(ctx context.Context, bucket, group string) error
}

type Service struct {
	Store
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		Store: NewStoreService(repo),
	}
}
