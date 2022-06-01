package grpc

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Alexander272/sealur/file_service/internal/models"
	proto_file "github.com/Alexander272/sealur/file_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/file_service/pkg/logger"
)

func (h *Handler) Download(req *proto_file.FileDownloadRequest, stream proto_file.FileService_DownloadServer) error {
	file, err := h.service.GetFile(context.Background(), req.Backet, req.Group, req.Id, req.Name)
	if err != nil {
		return fmt.Errorf("error getting file %w", err)
	}

	reqMeta := &proto_file.FileDownloadResponse{
		Response: &proto_file.FileDownloadResponse_Metadata{
			Metadata: &proto_file.MetaData{
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

		reqChunk := &proto_file.FileDownloadResponse{
			Response: &proto_file.FileDownloadResponse_File{File: &proto_file.File{
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

func (h *Handler) Upload(stream proto_file.FileService_UploadServer) error {
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

	id, err := h.service.Create(context.Background(), meta.Backet, fileDTO)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return stream.SendAndClose(&proto_file.FileUploadResponse{
		Id:       id,
		OrigName: meta.Name,
		Name:     fmt.Sprintf("%s_%s", id, meta.Name),
		Url:      fmt.Sprintf("/files/%s/%s/%s/%s", meta.Backet, meta.Group, id, meta.Name),
	})
}

func (h *Handler) Delete(ctx context.Context, req *proto_file.FileDeleteRequest) (*proto_file.FileDeleteResponse, error) {
	if err := h.service.Delete(ctx, req.Backet, req.Group, req.Id, req.Name); err != nil {
		return nil, err
	}
	return &proto_file.FileDeleteResponse{Message: "Removed"}, nil
}
