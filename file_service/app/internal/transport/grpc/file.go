package grpc

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Alexander272/sealur/file_service/internal/models"
	proto_file "github.com/Alexander272/sealur/file_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/file_service/pkg/logger"
)

//TODO надо запустить и посмотреть как именно там файл сохраняется
// func (h *Handler) Download(req *proto_file.FileDownloadRequest, stream proto_file.FileService_DownloadServer) error {
// 	file, err := h.service.GetFile(context.Background(), req.Uuid, )

// }

func (h *Handler) Upload(stream proto_file.FileService_UploadServer) error {
	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("cannot receive image info %w", err)
	}

	meta := req.GetMetadata()
	logger.Debug(meta.Backet)

	imageData := bytes.Buffer{}

	for {
		logger.Debug("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			logger.Debug("no more data")
			break
		}
		// if err == io.EOF {

		// 	file := req.GetFile().Content
		// 	data := req.GetMetadata()

		// 	reader := bytes.NewReader(file)

		// 	fileDTO := models.CreateFileDTO{
		// 		Name:   data.Name,
		// 		Size:   data.Size,
		// 		Reader: reader,
		// 	}

		// 	id, err := h.service.Create(context.Background(), data.Uuid, fileDTO)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to save file: %w", err)
		// 	}

		// 	return stream.SendAndClose(&proto_file.FileUploadResponse{
		// 		Id:   id,
		// 		Name: data.Name,
		// 	})
		// }
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
		Group:       meta.Uuid,
		ContentType: meta.Type,
		Reader:      reader,
	}

	id, err := h.service.Create(context.Background(), meta.Backet, fileDTO)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return stream.SendAndClose(&proto_file.FileUploadResponse{
		Id:   id,
		Name: meta.Name,
	})
}
