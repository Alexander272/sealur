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

func (s *StoreService) GetFile(ctx context.Context, bucket, group, id, name string) (f *models.File, err error) {
	f, err = s.repo.GetFile(ctx, bucket, fmt.Sprintf("%s/%s_%s", group, id, name))
	if err != nil {
		return f, err
	}
	return f, nil
}

func (s *StoreService) GetFilesByGroup(ctx context.Context, bucket, group string) ([]*models.File, error) {
	files, err := s.repo.GetFilesByGroup(ctx, bucket, group)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *StoreService) Create(ctx context.Context, bucket string, dto models.CreateFileDTO) (string, error) {
	dto.NormalizeName()
	file, err := models.NewFile(&dto)
	if err != nil {
		return "", err
	}
	err = s.repo.CreateFile(ctx, bucket, file)
	if err != nil {
		return "", err
	}
	return file.ID, nil
}

func (s *StoreService) Copy(ctx context.Context, bucket, group, newGroup, id string) error {
	err := s.repo.CopyFile(ctx, bucket, fmt.Sprintf("%s/%s", group, id), fmt.Sprintf("%s/%s", newGroup, id))
	if err != nil {
		return err
	}
	return nil
}

func (s *StoreService) CopyGroup(ctx context.Context, bucket, group, newGroup string) error {
	err := s.repo.CopyFiles(ctx, bucket, group, newGroup)
	if err != nil {
		return err
	}
	return nil
}

func (s *StoreService) Delete(ctx context.Context, bucket, group, id, name string) error {
	err := s.repo.DeleteFile(ctx, bucket, fmt.Sprintf("%s/%s_%s", group, id, name))
	if err != nil {
		return err
	}
	return nil
}

func (s *StoreService) DeleteGroup(ctx context.Context, bucket, group string) error {
	err := s.repo.DeleteFiles(ctx, bucket, group)
	if err != nil {
		return err
	}
	return nil
}
