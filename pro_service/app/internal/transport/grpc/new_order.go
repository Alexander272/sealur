package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/service"
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

func (h *OrderHandlers) Save(ctx context.Context, order order_api.CreateOrder) (*order_api.OrderNumber, error) {
	number, err := h.service.Save(ctx, order)
	if err != nil {
		return nil, err
	}
	return number, nil
}

func (h *OrderHandlers) Create(ctx context.Context, order order_api.CreateOrder) (*response_model.Response, error) {
	if err := h.service.Create(ctx, order); err != nil {
		return nil, err
	}
	return &response_model.Response{}, nil
}
