package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Alexander272/sealur/file_service/internal/models"
	"github.com/Alexander272/sealur/file_service/pkg/logger"
	"github.com/Alexander272/sealur/file_service/pkg/storage"
)

type StoreRepo struct {
	storage storage.Provider
}

func NewStoreRepo(storage storage.Provider) *StoreRepo {
	return &StoreRepo{storage: storage}
}

func (r *StoreRepo) GetFile(ctx context.Context, bucketName, fileName string) (*models.File, error) {
	// reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	obj, err := r.storage.GetFile(bucketName, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get file. err: %w", err)
	}
	defer obj.Close()

	objectInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info. err: %w", err)
	}

	buffer := make([]byte, objectInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}

	f := models.File{
		ID:          objectInfo.Key,
		ContentType: objectInfo.ContentType,
		Size:        objectInfo.Size,
		Bytes:       buffer,
	}
	return &f, nil

	// buffer := make([]byte, 1024)
	// _, err = obj.Read(buffer)
	// if err != nil && err != io.EOF {
	// 	return nil, fmt.Errorf("failed to get objects. err: %w", err)
	// }

	// f := models.File{
	// 	Bytes: buffer,
	// }
	// return &f, nil
}

func (r *StoreRepo) GetFilesByGroup(backet, group string) ([]*models.File, error) {
	objects, err := r.storage.GetBucketFiles(backet, group)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}
	if len(objects) == 0 {
		return nil, models.ErrNotFound
	}

	var files []*models.File
	for _, obj := range objects {
		stat, err := obj.Stat()
		if err != nil {
			logger.Errorf("failed to get objects. err: %v", err)
			continue
		}
		buffer := make([]byte, stat.Size)
		_, err = obj.Read(buffer)
		if err != nil && err != io.EOF {
			logger.Errorf("failed to get objects. err: %v", err)
			continue
		}
		f := models.File{
			ID:          stat.Key,
			Name:        strings.Split(stat.Key, "/")[1],
			ContentType: stat.ContentType,
			Size:        stat.Size,
			Bytes:       buffer,
		}
		files = append(files, &f)
		obj.Close()
	}

	return files, nil
}

func (r *StoreRepo) CreateFile(ctx context.Context, backet string, file *models.File) error {
	err := r.storage.UploadFile(
		fmt.Sprintf("%s/%s_%s", file.Group, file.ID, file.Name),
		file.Name,
		file.ContentType,
		backet,
		file.Size,
		bytes.NewBuffer(file.Bytes),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *StoreRepo) DeleteFile(ctx context.Context, backet, fileName string) error {
	err := r.storage.DeleteFile(backet, fileName)
	if err != nil {
		return err
	}
	return nil
}
