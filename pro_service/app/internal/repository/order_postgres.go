package repository

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) GetAll(req *proto.GetAllOrdersRequest) (orders []models.Order, err error) {
	query := fmt.Sprintf("SELECT id, date, count_position FROM %s WHERE user_id=$1", OrdersTable)

	if err = r.db.Select(&orders, query, req.UserId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return orders, nil
}

func (r *OrderRepo) Create(order *proto.CreateOrderRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (id, user_id, count_position) VALUES ($1, $2, $3)", OrdersTable)

	_, err := r.db.Exec(query, order.OrderId, order.UserId, order.Count)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Delete(order *proto.DeleteOrderRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", OrdersTable)

	_, err := r.db.Exec(query, order.OrderId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Save(order *proto.SaveOrderRequest) error {
	// TODO дописать сохранение
	return nil
}
