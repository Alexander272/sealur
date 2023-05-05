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

func (h *OrderHandlers) Get(ctx context.Context, req *order_api.GetOrder) (*order_api.Order, error) {
	order, err := h.service.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return &order_api.Order{Order: order}, nil
}

func (h *OrderHandlers) GetCurrent(ctx context.Context, req *order_api.GetCurrentOrder) (*order_model.CurrentOrder, error) {
	order, err := h.service.GetCurrent(ctx, req)
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

func (h *OrderHandlers) GetOpen(ctx context.Context, req *order_api.GetManagerOrders) (*order_api.ManagerOrders, error) {
	orders, err := h.service.GetOpen(ctx, req)
	if err != nil {
		return nil, err
	}
	return &order_api.ManagerOrders{Orders: orders}, nil
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

func (h *OrderHandlers) GetAnalytics(ctx context.Context, req *order_api.GetOrderAnalytics) (*order_api.Analytics, error) {
	analytics, err := h.service.GetAnalytics(ctx, req)
	if err != nil {
		return nil, err
	}
	return analytics, nil
}

func (h *OrderHandlers) GetBidAnalytics(ctx context.Context, req *order_api.GetFullOrderAnalytics) (*order_api.OrderAnalytics, error) {
	analytics, err := h.service.GetFullAnalytics(ctx, req)
	if err != nil {
		return nil, err
	}
	return &order_api.OrderAnalytics{Orders: analytics}, nil
}

func (h *OrderHandlers) Save(ctx context.Context, order *order_api.CreateOrder) (*order_api.OrderNumber, error) {
	number, err := h.service.Save(ctx, order)
	if err != nil {
		return nil, err
	}
	return number, nil
}

func (h *OrderHandlers) Copy(ctx context.Context, order *order_api.CopyOrder) (*response_model.Response, error) {
	if err := h.service.Copy(ctx, order); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *OrderHandlers) Create(ctx context.Context, order *order_api.CreateOrder) (*response_model.IdResponse, error) {
	id, err := h.service.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	return &response_model.IdResponse{Id: id}, nil
}

func (h *OrderHandlers) SetStatus(ctx context.Context, status *order_api.Status) (*response_model.Response, error) {
	if err := h.service.SetStatus(ctx, status); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *OrderHandlers) SetInfo(ctx context.Context, info *order_api.Info) (*response_model.Response, error) {
	if err := h.service.SetInfo(ctx, info); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}

func (h *OrderHandlers) SetManager(ctx context.Context, manager *order_api.Manager) (*response_model.Response, error) {
	if err := h.service.SetManager(ctx, manager); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
