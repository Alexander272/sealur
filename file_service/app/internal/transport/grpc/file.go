package grpc

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Alexander272/sealur/file_service/internal/models"
	"github.com/Alexander272/sealur/file_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/file_api"
)

func (h *Handler) Download(req *file_api.FileDownloadRequest, stream file_api.FileService_DownloadServer) error {
	file, err := h.service.GetFile(context.Background(), req.Bucket, req.Group, req.Id, req.Name)
	if err != nil {
		return fmt.Errorf("error getting file %w", err)
	}

	reqMeta := &file_api.FileDownloadResponse{
		Response: &file_api.FileDownloadResponse_Metadata{
			Metadata: &file_api.MetaData{
				Name:  file.Name,
				Size:  file.Size,
				Type:  file.ContentType,
				Group: file.Group,
			},
		},
	}
	err = stream.Send(reqMeta)
	if err != nil {
		logger.Errorf("cannot send metadata to clinet: %w", err)
		return fmt.Errorf("cannot send metadata to clinet %w", err)
	}

	f := bytes.NewReader(file.Bytes)
	reader := bufio.NewReader(f)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf("cannot read chunk to buffer: %w", err)
			return fmt.Errorf("cannot read chunk to buffer %w", err)
		}

		reqChunk := &file_api.FileDownloadResponse{
			Response: &file_api.FileDownloadResponse_File{File: &file_api.File{
				Content: buffer[:n],
			}},
		}

		err = stream.Send(reqChunk)
		if err != nil {
			logger.Errorf("cannot send chunk to clinet: %w", err)
			return fmt.Errorf("cannot send chunk to clinet %w", err)
		}
	}

	return nil
}

func (h *Handler) GroupDownload(req *file_api.GroupDownloadRequest, stream file_api.FileService_GroupDownloadServer) error {
	files, err := h.service.GetFilesByGroup(context.Background(), req.Bucket, req.Group)
	if err != nil {
		return fmt.Errorf("error getting files. error: %w", err)
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return fmt.Errorf("failed to create file in zip. err %w", err)
		}
		_, err = f.Write(file.Bytes)
		if err != nil {
			return fmt.Errorf("failed to write file in zip. err %w", err)
		}
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close writer. err %w", err)
	}

	size := int64(buf.Cap())
	reqMeta := &file_api.FileDownloadResponse{
		Response: &file_api.FileDownloadResponse_Metadata{
			Metadata: &file_api.MetaData{
				Name: "Чертежи.zip",
				Size: size,
				Type: ".zip",
			},
		},
	}
	err = stream.Send(reqMeta)
	if err != nil {
		logger.Errorf("cannot send metadata to clinet: %w", err)
		return fmt.Errorf("cannot send metadata to clinet %w", err)
	}

	reader := bufio.NewReader(buf)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf("cannot read chunk to buffer: %w", err)
			return fmt.Errorf("cannot read chunk to buffer %w", err)
		}

		reqChunk := &file_api.FileDownloadResponse{
			Response: &file_api.FileDownloadResponse_File{File: &file_api.File{
				Content: buffer[:n],
			}},
		}

		err = stream.Send(reqChunk)
		if err != nil {
			logger.Errorf("cannot send chunk to clinet: %w", err)
			return fmt.Errorf("cannot send chunk to clinet %w", err)
		}
	}

	return nil
}

func (h *Handler) Upload(stream file_api.FileService_UploadServer) error {
	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("cannot receive image info %w", err)
	}

	meta := req.GetMetadata()
	imageData := bytes.Buffer{}

	for {
		logger.Debug("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			logger.Debug("no more data")
			break
		}

		if err != nil {
			return fmt.Errorf("cannot receive chunk data: %w", err)
		}

		chunk := req.GetFile().Content

		_, err = imageData.Write(chunk)
		if err != nil {
			return fmt.Errorf("cannot write chunk data: %w", err)
		}
	}

	reader := bytes.NewReader(imageData.Bytes())

	fileDTO := models.CreateFileDTO{
		Name:        meta.Name,
		Size:        meta.Size,
		Group:       meta.Group,
		ContentType: meta.Type,
		Reader:      reader,
	}

	id, err := h.service.Create(context.Background(), meta.Bucket, fileDTO)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return stream.SendAndClose(&file_api.FileUploadResponse{
		Id:       id,
		OrigName: meta.Name,
		Name:     fmt.Sprintf("%s_%s", id, meta.Name),
		Url:      fmt.Sprintf("/files/%s/%s/%s/%s", meta.Bucket, meta.Group, id, meta.Name),
	})
}

func (h *Handler) Copy(ctx context.Context, req *file_api.CopyFileRequest) (*file_api.MessageResponse, error) {
	if err := h.service.Copy(ctx, req.Bucket, req.Group, req.NewGroup, req.Id); err != nil {
		return nil, err
	}
	return &file_api.MessageResponse{Message: "Copied"}, nil
}

func (h *Handler) CopyGroup(ctx context.Context, req *file_api.CopyGroupRequest) (*file_api.MessageResponse, error) {
	if err := h.service.CopyGroup(ctx, req.Bucket, req.Group, req.NewGroup); err != nil {
		return nil, err
	}
	return &file_api.MessageResponse{Message: "Copied"}, nil
}

func (h *Handler) Delete(ctx context.Context, req *file_api.FileDeleteRequest) (*file_api.MessageResponse, error) {
	if err := h.service.Delete(ctx, req.Bucket, req.Group, req.Id, req.Name); err != nil {
		return nil, err
	}
	return &file_api.MessageResponse{Message: "Removed"}, nil
}

func (h *Handler) GroupDelete(ctx context.Context, req *file_api.GroupDeleteRequest) (*file_api.MessageResponse, error) {
	if err := h.service.DeleteGroup(ctx, req.Bucket, req.Group); err != nil {
		return nil, err
	}
	return &file_api.MessageResponse{Message: "Removed"}, nil
}
