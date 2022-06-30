package service

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	proto_email "github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto/email"
	proto_file "github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto/file"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/xuri/excelize/v2"
)

type OrderService struct {
	repo  repository.Order
	email proto_email.EmailServiceClient
	file  proto_file.FileServiceClient
}

var columnNames = []interface{}{"№", "Обозначение", "Количество", "Описание", "Размеры", "Чертеж"}

func NewOrderService(repo repository.Order, email proto_email.EmailServiceClient, file proto_file.FileServiceClient) *OrderService {
	return &OrderService{
		repo:  repo,
		email: email,
		file:  file,
	}
}

func (s *OrderService) GetAll(req *proto.GetAllOrdersRequest) (orders []*proto.Order, err error) {
	o, err := s.repo.GetAll(req)
	if err != nil {
		return orders, fmt.Errorf("failed to get orders. error: %w", err)
	}

	for _, item := range o {
		order := proto.Order(item)
		orders = append(orders, &order)
	}

	return orders, nil
}

func (s *OrderService) Create(order *proto.CreateOrderRequest) (*proto.IdResponse, error) {
	o, err := s.repo.GetCur(&proto.GetCurOrderRequest{UserId: order.UserId})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get order. error: %w", err)
	}
	if (o != models.Order{}) {
		return &proto.IdResponse{Id: o.Id}, nil
	}

	if err := s.repo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order. error: %w", err)
	}
	return &proto.IdResponse{Id: order.OrderId}, nil
}

func (s *OrderService) Delete(order *proto.DeleteOrderRequest) (*proto.IdResponse, error) {
	_, err := s.file.GroupDelete(context.Background(), &proto_file.GroupDeleteRequest{
		Bucket: "pro",
		Group:  order.OrderId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete group files. error: %w", err)
	}

	if err := s.repo.Delete(order); err != nil {
		return nil, fmt.Errorf("failed to delete order. error: %w", err)
	}
	return &proto.IdResponse{Id: order.OrderId}, nil
}

func (s *OrderService) Copy(order *proto.CopyOrderRequest) error {
	_, err := s.file.CopyGroup(context.Background(), &proto_file.CopyGroupRequest{
		Bucket:   "pro",
		Group:    order.OldOrderId,
		NewGroup: order.OrderId,
	})
	if err != nil {
		return fmt.Errorf("failed to copy group files. error: %w", err)
	}

	if err := s.repo.Copy(order); err != nil {
		return fmt.Errorf("failed to copy order. error: %w", err)
	}

	return nil
}

func (s *OrderService) Save(ctx context.Context, order *proto.SaveOrderRequest) (*bytes.Buffer, error) {
	return s.createZip(ctx, order)
}

func (s *OrderService) Send(ctx context.Context, order *proto.SaveOrderRequest) error {
	_, err := s.createZip(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) createZip(ctx context.Context, order *proto.SaveOrderRequest) (*bytes.Buffer, error) {
	if err := s.repo.Save(order); err != nil {
		return nil, fmt.Errorf("failed to save order. error: %w", err)
	}

	positions, err := s.repo.GetPositions(&proto.GetPositionsRequest{OrderId: order.OrderId})
	if err != nil {
		return nil, fmt.Errorf("failed to get positions. error: %w", err)
	}

	file := excelize.NewFile()
	sheetName := file.GetSheetName(file.GetActiveSheetIndex())

	streamWriter, err := file.NewStreamWriter(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream writer. error: %w", err)
	}

	cell, _ := excelize.CoordinatesToCellName(1, 1)
	if err := streamWriter.SetRow(cell, columnNames); err != nil {
		return nil, fmt.Errorf("failed to create header table. error: %w", err)
	}

	for i, p := range positions {
		line := []interface{}{i + 1, p.Designation, p.Count, p.Descriprion, p.Sizes, p.Drawing}

		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := streamWriter.SetRow(cell, line); err != nil {
			return nil, fmt.Errorf("failed to create line table. error: %w", err)
		}
	}
	if err := streamWriter.Flush(); err != nil {
		return nil, fmt.Errorf("failed to close stream writer. error: %w", err)
	}

	stream, err := s.file.GroupDownload(ctx, &proto_file.GroupDownloadRequest{
		Bucket: "pro",
		Group:  order.OrderId,
	})
	if err != nil {
		logger.Errorf("failed to download drawing. err :%w", err)
		return nil, fmt.Errorf("failed to download drawing. err :%w", err)
	}

	req, err := stream.Recv()
	if err != nil && !strings.Contains(err.Error(), "file not found") {
		return nil, fmt.Errorf("failed to get data. err: %w", err)
	}
	meta := req.GetMetadata()
	fileData := bytes.Buffer{}

	if meta != nil {
		for {
			logger.Debug("waiting to receive more data")

			req, err := stream.Recv()
			if err == io.EOF {
				logger.Debug("no more data")
				break
			}

			if err != nil {
				return nil, fmt.Errorf("failed to get chunk. err %w", err)
			}

			chunk := req.GetFile().Content

			_, err = fileData.Write(chunk)
			if err != nil {
				return nil, fmt.Errorf("failed to write chunk. err %w", err)
			}
		}
	}

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	fw, err := writer.Create("Заказ.xlsx")
	if err != nil {
		return nil, fmt.Errorf("failed to create xlsx in zip. err %w", err)
	}
	_, err = file.WriteTo(fw)
	if err != nil {
		return nil, fmt.Errorf("failed to write xlsx in zip. err %w", err)
	}

	if meta != nil {
		fw, err := writer.Create(meta.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to create zip in zip. err %w", err)
		}
		_, err = fw.Write(fileData.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to write zip in zip. err %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer. err %w", err)
	}

	return buf, nil
}
