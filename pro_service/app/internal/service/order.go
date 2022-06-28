package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
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
	if err := s.repo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order. error: %w", err)
	}
	return &proto.IdResponse{Id: order.OrderId}, nil
}

func (s *OrderService) Delete(order *proto.DeleteOrderRequest) (*proto.IdResponse, error) {
	if err := s.repo.Delete(order); err != nil {
		return nil, fmt.Errorf("failed to delete order. error: %w", err)
	}
	return &proto.IdResponse{Id: order.OrderId}, nil
}

func (s *OrderService) Save(order *proto.SaveOrderRequest) error {
	if err := s.repo.Save(order); err != nil {
		return fmt.Errorf("failed to save order. error: %w", err)
	}
	return nil
}
