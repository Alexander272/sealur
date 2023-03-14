package grpc

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/Alexander272/sealur/pro_service/internal/service"
	"github.com/Alexander272/sealur_proto/api/pro/models/order_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/response_model"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
)

type OrderHandlers struct {
	service service.OrderNew
	order_api.UnimplementedOrderServiceServer
}

func NewOrderHandlers(service service.OrderNew) *OrderHandlers {
	return &OrderHandlers{
		service: service,
	}
}

func (h *OrderHandlers) Get(ctx context.Context, req *order_api.GetOrder) (*order_model.FullOrder, error) {
	order, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (h *OrderHandlers) GetAll(ctx context.Context, req *order_api.GetAllOrders) (*order_api.Orders, error) {
	orders, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return &order_api.Orders{Orders: orders}, nil
}

func (h *OrderHandlers) GetFile(req *order_api.GetOrder, stream order_api.OrderService_GetFileServer) error {
	file, fileName, err := h.service.GetFile(context.Background(), req)
	if err != nil {
		return err
	}

	reqMeta := &response_model.FileResponse{
		Response: &response_model.FileResponse_Metadata{
			Metadata: &response_model.MetaData{
				Name: fileName + ".zip",
				Size: int64(file.Cap()),
				Type: "application/zip",
			},
		},
	}
	err = stream.Send(reqMeta)
	if err != nil {
		return fmt.Errorf("cannot send metadata to client %w", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("cannot read chunk to buffer %w", err)
		}

		reqChunk := &response_model.FileResponse{
			Response: &response_model.FileResponse_File{File: &response_model.File{
				Content: buffer[:n],
			}},
		}

		err = stream.Send(reqChunk)
		if err != nil {
			return fmt.Errorf("cannot send chunk to client %w", err)
		}
	}

	return nil
}

func (h *OrderHandlers) Save(ctx context.Context, order *order_api.CreateOrder) (*order_api.OrderNumber, error) {
	number, err := h.service.Save(ctx, order)
	if err != nil {
		return nil, err
	}
	return number, nil
}

func (h *OrderHandlers) Create(ctx context.Context, order *order_api.CreateOrder) (*response_model.Response, error) {
	if err := h.service.Create(ctx, order); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
