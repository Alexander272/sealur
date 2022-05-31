package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Alexander272/sealur/file_service/internal/config"
	"github.com/Alexander272/sealur/file_service/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client *minio.Client
}

func NewClient(conf config.MinIOConfig) (*MinioStorage, error) {
	logger.Debug(conf.Endpoint)
	minioClient, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}

	return &MinioStorage{Client: minioClient}, nil
}

func (c *MinioStorage) GetFile(ctx context.Context, bucketName, fileId string) (*minio.Object, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	obj, err := c.Client.GetObject(reqCtx, bucketName, fileId, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file with id: %s from minio bucket %s. err: %w", fileId, bucketName, err)
	}
	return obj, nil
}

func (c *MinioStorage) GetBucketFiles(ctx context.Context, bucketName string) ([]*minio.Object, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var files []*minio.Object
	for lobj := range c.Client.ListObjects(reqCtx, bucketName, minio.ListObjectsOptions{WithMetadata: true}) {
		if lobj.Err != nil {
			logger.Errorf("failed to list object from minio bucket %s. err: %v", bucketName, lobj.Err)
			continue
		}
		object, err := c.Client.GetObject(ctx, bucketName, lobj.Key, minio.GetObjectOptions{})
		if err != nil {
			logger.Errorf("failed to get object key=%s from minio bucket %s. err: %v", lobj.Key, bucketName, lobj.Err)
			continue
		}
		files = append(files, object)
	}
	return files, nil
}

func (c *MinioStorage) UploadFile(ctx context.Context, fileId, fileName, contetnType, bucketName string, fileSize int64, reader io.Reader) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	exists, errBucketExists := c.Client.BucketExists(ctx, bucketName)
	if errBucketExists != nil || !exists {
		logger.Warnf("no bucket %s. creating new one...", bucketName)
		err := c.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create new bucket. err: %w", err)
		}
	}

	logger.Debugf("put new object %s to bucket %s", fileName, bucketName)
	_, err := c.Client.PutObject(
		reqCtx,
		bucketName,
		fileId,
		reader,
		fileSize,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": fileName,
			},
			ContentType: contetnType,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload file. err: %w", err)
	}
	return nil
}

func (c *MinioStorage) DeleteFile(ctx context.Context, noteUUID, fileName string) error {
	err := c.Client.RemoveObject(ctx, noteUUID, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file. err: %w", err)
	}
	return nil
}