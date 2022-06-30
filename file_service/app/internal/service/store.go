package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/file_service/internal/models"
	"github.com/Alexander272/sealur/file_service/internal/repository"
)

type StoreService struct {
	repo *repository.Repo
}

func NewStoreService(repo *repository.Repo) *StoreService {
	return &StoreService{repo: repo}
}

func (s *StoreService) GetFile(ctx context.Context, backet, group, id, name string) (f *models.File, err error) {
	f, err = s.repo.GetFile(ctx, backet, fmt.Sprintf("%s/%s_%s", group, id, name))
	if err != nil {
		return f, err
	}
	return f, nil
}

func (s *StoreService) GetFilesByGroup(ctx context.Context, backet, group string) ([]*models.File, error) {
	files, err := s.repo.GetFilesByGroup(backet, group)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *StoreService) Create(ctx context.Context, backet string, dto models.CreateFileDTO) (string, error) {
	dto.NormalizeName()
	file, err := models.NewFile(&dto)
	if err != nil {
		return "", err
	}
	err = s.repo.CreateFile(ctx, backet, file)
	if err != nil {
		return "", err
	}
	return file.ID, nil
}

func (s *StoreService) Delete(ctx context.Context, backet, group, id, name string) error {
	err := s.repo.DeleteFile(ctx, backet, fmt.Sprintf("%s/%s_%s", group, id, name))
	if err != nil {
		return err
	}
	return nil
}
