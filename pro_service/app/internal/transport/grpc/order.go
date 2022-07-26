package grpc

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetAllOrders(ctx context.Context, req *pro_api.GetAllOrdersRequest) (*pro_api.OrderResponse, error) {
	orders, err := h.service.Order.GetAll(req)
	if err != nil {
		return nil, err
	}

	return &pro_api.OrderResponse{Orders: orders}, nil
}

func (h *Handler) CreateOrder(ctx context.Context, order *pro_api.CreateOrderRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.Order.Create(order)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) DeleteOrder(ctx context.Context, order *pro_api.DeleteOrderRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.Order.Delete(order)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) SaveOrder(order *pro_api.SaveOrderRequest, stream pro_api.ProService_SaveOrderServer) error {
	file, err := h.service.Order.Save(context.Background(), order)
	if err != nil {
		return err
	}

	reqMeta := &pro_api.FileDownloadResponse{
		Response: &pro_api.FileDownloadResponse_Metadata{
			Metadata: &pro_api.MetaData{
				Name: "Order.zip",
				Size: int64(file.Cap()),
				Type: "application/zip",
			},
		},
	}
	err = stream.Send(reqMeta)
	if err != nil {
		return fmt.Errorf("cannot send metadata to clinet %w", err)
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

		reqChunk := &pro_api.FileDownloadResponse{
			Response: &pro_api.FileDownloadResponse_File{File: &pro_api.File{
				Content: buffer[:n],
			}},
		}

		err = stream.Send(reqChunk)
		if err != nil {
			return fmt.Errorf("cannot send chunk to clinet %w", err)
		}
	}

	return nil
}

func (h *Handler) SendOrder(ctx context.Context, order *pro_api.SaveOrderRequest) (*pro_api.SuccessResponse, error) {
	if err := h.service.Order.Send(ctx, order); err != nil {
		return nil, err
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (h *Handler) CopyOrder(ctx context.Context, order *pro_api.CopyOrderRequest) (*pro_api.SuccessResponse, error) {
	if err := h.service.Order.Copy(order); err != nil {
		return nil, err
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (h *Handler) GetPositions(ctx context.Context, req *pro_api.GetPositionsRequest) (*pro_api.OrderPositionResponse, error) {
	positions, err := h.service.OrderPosition.Get(req)
	if err != nil {
		return nil, err
	}

	return &pro_api.OrderPositionResponse{Positions: positions}, nil
}

func (h *Handler) GetCurPositions(ctx context.Context, req *pro_api.GetCurPositionsRequest) (*pro_api.OrderPositionResponse, error) {
	positions, err := h.service.OrderPosition.GetCur(req)
	if err != nil {
		return nil, err
	}

	return &pro_api.OrderPositionResponse{Positions: positions}, nil
}

func (h *Handler) AddPosition(ctx context.Context, position *pro_api.AddPositionRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.OrderPosition.Add(position)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) CopyPosition(ctx context.Context, position *pro_api.CopyPositionRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.OrderPosition.Copy(position)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) UpdatePosition(ctx context.Context, position *pro_api.UpdatePositionRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.OrderPosition.Update(position)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) RemovePosition(ctx context.Context, position *pro_api.RemovePositionRequest) (*pro_api.IdResponse, error) {
	id, err := h.service.OrderPosition.Remove(position)
	if err != nil {
		return nil, err
	}
	return id, nil
}
