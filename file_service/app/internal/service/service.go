package service

import (
	"context"

	"github.com/Alexander272/sealur/file_service/internal/models"
	"github.com/Alexander272/sealur/file_service/internal/repository"
)

type Store interface {
	GetFile(ctx context.Context, orderUUID, fileName string) (*models.File, error)
	GetFilesByOrderUUID(ctx context.Context, orderUUID string) ([]*models.File, error)
	Create(ctx context.Context, backet string, dto models.CreateFileDTO) (string, error)
	Delete(ctx context.Context, backet, fileName string) error
}

type Service struct {
	Store
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		Store: NewStoreService(repo),
	}
}
