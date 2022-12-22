package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type File struct {
	Name string
	Url  string
}

type Provider interface {
	GetFile(ctx context.Context, bucketName, fileId string) (*minio.Object, error)
	GetBucketFiles(ctx context.Context, bucketName, group string) ([]*minio.Object, error)
	UploadFile(ctx context.Context, fileId, fileName, contetnType, bucketName string, fileSize int64, reader io.Reader) error
	CopyFile(ctx context.Context, destBucket, desctFileId, srcBucket, srcFileId string) error
	CopyGroupFiles(ctx context.Context, bucket, group, newGroup string) error
	DeleteFile(ctx context.Context, bucket, fileId string) error
	DeleteGroupFiles(ctx context.Context, bucket, group string) error
}
