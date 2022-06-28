package grpc

import (
	"context"

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

func (h *Handler) SaveOrder(ctx context.Context, order *proto.SaveOrderRequest) (*proto.SuccessResponse, error) {
	if err := h.service.Order.Save(order); err != nil {
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
