package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro/models/order_model"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/google/uuid"
)

type OrderServiceNew struct {
	repo     repository.OrderNew
	position Position
}

func NewOrderService_New(repo repository.OrderNew, position Position) *OrderServiceNew {
	return &OrderServiceNew{
		repo:     repo,
		position: position,
	}
}

func (s *OrderServiceNew) Get(ctx context.Context, req *order_api.GetOrder) (*order_model.FullOrder, error) {
	order, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id. error: %w", err)
	}

	// TODO получить список позиций с данными

	return order, nil
}

func (s *OrderServiceNew) GetAll(ctx context.Context, req *order_api.GetAllOrders) ([]*order_model.Order, error) {
	orders, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders. error: %w", err)
	}
	return orders, nil
}

func (s *OrderServiceNew) Save(ctx context.Context, order *order_api.CreateOrder) (*order_api.OrderNumber, error) {
	var orderId = order.Id
	var number int64
	if orderId == "" {
		orderId = uuid.New().String()

		order.Id = orderId
		if err := s.repo.Create(ctx, order, fmt.Sprintf("%d", time.Now().UnixMilli())); err != nil {
			return nil, fmt.Errorf("failed to create order. error: %w", err)
		}
	} else {
		var err error
		number, err = s.repo.GetNumber(ctx, orderId, fmt.Sprintf("%d", time.Now().UnixMilli()))
		if err != nil {
			return nil, fmt.Errorf("failed to get order number. error: %w", err)
		}
	}

	if err := s.position.CreateSeveral(ctx, order.Positions, orderId); err != nil {
		return nil, err
	}

	return &order_api.OrderNumber{Number: number}, nil
}

func (s *OrderServiceNew) Create(ctx context.Context, order *order_api.CreateOrder) error {
	var orderId = order.Id
	if orderId == "" {
		orderId = uuid.New().String()

		order.Id = orderId
		if err := s.repo.Create(ctx, order, ""); err != nil {
			return fmt.Errorf("failed to create order. error: %w", err)
		}
	}

	if err := s.position.CreateSeveral(ctx, order.Positions, orderId); err != nil {
		return err
	}

	return nil
}
