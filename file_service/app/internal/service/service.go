package service

import (
	"context"

	"github.com/Alexander272/sealur/file_service/internal/models"
	"github.com/Alexander272/sealur/file_service/internal/repository"
)

type Store interface {
	GetFile(ctx context.Context, backet, group, id, name string) (*models.File, error)
	// GetFilesByOrderUUID(ctx context.Context, backet string) ([]*models.File, error)
	Create(ctx context.Context, backet string, dto models.CreateFileDTO) (string, error)
	Delete(ctx context.Context, backet, group, id, name string) error
}

type Service struct {
	Store
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		Store: NewStoreService(repo),
	}
}
