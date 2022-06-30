package repository

import (
	"context"

	"github.com/Alexander272/sealur/file_service/internal/models"
	"github.com/Alexander272/sealur/file_service/pkg/storage"
)

type Store interface {
	GetFile(ctx context.Context, bucketName, fileName string) (*models.File, error)
	GetFilesByGroup(ctx context.Context, bucketName, group string) ([]*models.File, error)
	CreateFile(ctx context.Context, bucket string, file *models.File) error
	CopyFile(ctx context.Context, bucket, fileName, newFileName string) error
	CopyFiles(ctx context.Context, bucket, group, newGroup string) error
	DeleteFile(ctx context.Context, bucket, fileName string) error
	DeleteFiles(ctx context.Context, bucket, group string) error
}

type Repo struct {
	Store
}

func NewRepo(storage storage.Provider) *Repo {
	return &Repo{
		Store: NewStoreRepo(storage),
	}
}
