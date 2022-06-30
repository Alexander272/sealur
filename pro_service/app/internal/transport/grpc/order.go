package grpc

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetAllOrders(ctx context.Context, req *proto.GetAllOrdersRequest) (*proto.OrderResponse, error) {
	orders, err := h.service.Order.GetAll(req)
	if err != nil {
		return nil, err
	}

	return &proto.OrderResponse{Orders: orders}, nil
}

func (h *Handler) CreateOrder(ctx context.Context, order *proto.CreateOrderRequest) (*proto.IdResponse, error) {
	id, err := h.service.Order.Create(order)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) DeleteOrder(ctx context.Context, order *proto.DeleteOrderRequest) (*proto.IdResponse, error) {
	id, err := h.service.Order.Delete(order)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) SaveOrder(order *proto.SaveOrderRequest, stream proto.ProService_SaveOrderServer) error {
	file, err := h.service.Order.Save(context.Background(), order)
	if err != nil {
		return err
	}

	reqMeta := &proto.FileDownloadResponse{
		Response: &proto.FileDownloadResponse_Metadata{
			Metadata: &proto.MetaData{
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

		reqChunk := &proto.FileDownloadResponse{
			Response: &proto.FileDownloadResponse_File{File: &proto.File{
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

func (h *Handler) SendOrder(ctx context.Context, order *proto.SaveOrderRequest) (*proto.SuccessResponse, error) {
	if err := h.service.Order.Send(ctx, order); err != nil {
		return nil, err
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (h *Handler) CopyOrder(ctx context.Context, order *proto.CopyOrderRequest) (*proto.SuccessResponse, error) {
	if err := h.service.Order.Copy(order); err != nil {
		return nil, err
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (h *Handler) GetPositions(ctx context.Context, req *proto.GetPositionsRequest) (*proto.OrderPositionResponse, error) {
	positions, err := h.service.OrderPosition.Get(req)
	if err != nil {
		return nil, err
	}

	return &proto.OrderPositionResponse{Positions: positions}, nil
}

func (h *Handler) GetCurPositions(ctx context.Context, req *proto.GetCurPositionsRequest) (*proto.OrderPositionResponse, error) {
	positions, err := h.service.OrderPosition.GetCur(req)
	if err != nil {
		return nil, err
	}

	return &proto.OrderPositionResponse{Positions: positions}, nil
}

func (h *Handler) AddPosition(ctx context.Context, position *proto.AddPositionRequest) (*proto.IdResponse, error) {
	id, err := h.service.OrderPosition.Add(position)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) CopyPosition(ctx context.Context, position *proto.CopyPositionRequest) (*proto.IdResponse, error) {
	id, err := h.service.OrderPosition.Copy(position)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) UpdatePosition(ctx context.Context, position *proto.UpdatePositionRequest) (*proto.IdResponse, error) {
	id, err := h.service.OrderPosition.Update(position)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (h *Handler) RemovePosition(ctx context.Context, position *proto.RemovePositionRequest) (*proto.IdResponse, error) {
	id, err := h.service.OrderPosition.Remove(position)
	if err != nil {
		return nil, err
	}
	return id, nil
}
