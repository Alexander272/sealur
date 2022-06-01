package storage

import (
	"io"

	"github.com/minio/minio-go"
)

type File struct {
	Name string
	Url  string
}

type Provider interface {
	GetFile(bucketName, fileId string) (*minio.Object, error)
	// GetBucketFiles(ctx context.Context, bucketName string) ([]*minio.Object, error)
	UploadFile(fileId, fileName, contetnType, bucketName string, fileSize int64, reader io.Reader) error
	DeleteFile(noteUUID, fileId string) error
}
