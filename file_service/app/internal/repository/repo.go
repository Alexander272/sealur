package repository

import (
	"context"

	"github.com/Alexander272/sealur/file_service/internal/models"
	"github.com/Alexander272/sealur/file_service/pkg/storage"
)

type Store interface {
	GetFile(ctx context.Context, bucketName, fileName string) (*models.File, error)
	GetFilesByOrderUUID(ctx context.Context, bucketName string) ([]*models.File, error)
	CreateFile(ctx context.Context, backet string, file *models.File) error
	DeleteFile(ctx context.Context, backet, fileName string) error
}

type Repo struct {
	Store
}

func NewRepo(storage storage.Provider) *Repo {
	return &Repo{
		Store: NewStoreRepo(storage),
	}
}
