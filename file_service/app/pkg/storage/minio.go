package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Alexander272/sealur/file_service/internal/config"
	"github.com/Alexander272/sealur/file_service/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client *minio.Client
}

func NewClient(conf config.MinIOConfig) (*MinioStorage, error) {
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
	obj, err := c.Client.GetObject(ctx, bucketName, fileId, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file with id: %s from minio bucket %s. err: %w", fileId, bucketName, err)
	}
	return obj, nil
}

func (c *MinioStorage) GetBucketFiles(ctx context.Context, bucketName, group string) ([]*minio.Object, error) {
	var files []*minio.Object

	for lobj := range c.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Prefix: group, Recursive: true}) {
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
		ctx,
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

func (c *MinioStorage) CopyFile(ctx context.Context, destBucket, desctFileId, srcBucket, srcFileId string) error {
	// Source object
	srcOpts := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcFileId,
	}

	// Destination object
	dstOpts := minio.CopyDestOptions{
		Bucket: destBucket,
		Object: desctFileId,
	}

	if _, err := c.Client.CopyObject(ctx, dstOpts, srcOpts); err != nil {
		return fmt.Errorf("failed to copy file. error: %w", err)
	}
	return nil
}

func (c *MinioStorage) CopyGroupFiles(ctx context.Context, bucket, group, newGroup string) error {
	for lobj := range c.Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: group, Recursive: true}) {
		if lobj.Err != nil {
			logger.Errorf("failed to list object from minio bucket %s. err: %v", bucket, lobj.Err)
			continue
		}

		// Source object
		srcOpts := minio.CopySrcOptions{
			Bucket: bucket,
			Object: lobj.Key,
		}

		// Destination object
		dstOpts := minio.CopyDestOptions{
			Bucket: bucket,
			Object: fmt.Sprintf("%s/%s", newGroup, strings.Split(lobj.Key, "/")[1]),
		}

		if _, err := c.Client.CopyObject(ctx, dstOpts, srcOpts); err != nil {
			return fmt.Errorf("failed to copy file. error: %w", err)
		}
	}

	return nil
}

func (c *MinioStorage) DeleteFile(ctx context.Context, bucket, fileId string) error {
	err := c.Client.RemoveObject(ctx, bucket, fileId, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file. err: %w", err)
	}
	return nil
}

func (c *MinioStorage) DeleteGroupFiles(ctx context.Context, bucket, group string) error {
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)

		for lobj := range c.Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: group, Recursive: true}) {
			if lobj.Err != nil {
				logger.Errorf("failed to list object from minio bucket %s. err: %v", bucket, lobj.Err)
				continue
			}
			objectsCh <- lobj
		}
	}()

	errorCh := c.Client.RemoveObjects(ctx, bucket, objectsCh, minio.RemoveObjectsOptions{})
	for e := range errorCh {
		return fmt.Errorf("failed to remove " + e.ObjectName + ", error: " + e.Err.Error())
	}
	return nil
}
