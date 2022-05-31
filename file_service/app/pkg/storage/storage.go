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
	GetBucketFiles(ctx context.Context, bucketName string) ([]*minio.Object, error)
	UploadFile(ctx context.Context, fileId, fileName, contetnType, bucketName string, fileSize int64, reader io.Reader) error
	DeleteFile(ctx context.Context, noteUUID, fileId string) error
}
