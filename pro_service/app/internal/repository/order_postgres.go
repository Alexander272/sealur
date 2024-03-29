package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) GetAll(req *pro_api.GetAllOrdersRequest) (orders []models.Order, err error) {
	query := fmt.Sprintf("SELECT id, date, count_position FROM %s WHERE user_id=$1 AND date IS NOT NULL", OrdersTable)

	if err = r.db.Select(&orders, query, req.UserId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return orders, nil
}

func (r *OrderRepo) GetCur(req *pro_api.GetCurOrderRequest) (order models.Order, err error) {
	query := fmt.Sprintf("SELECT id, count_position FROM %s WHERE user_id=$1 AND date IS NULL", OrdersTable)

	if err = r.db.Get(&order, query, req.UserId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order, err
		}
		return order, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return order, nil
}

func (r *OrderRepo) Create(order *pro_api.CreateOrderRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (id, user_id, count_position) VALUES ($1, $2, $3)", OrdersTable)

	_, err := r.db.Exec(query, order.OrderId, order.UserId, order.Count)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Copy(order *pro_api.CopyOrderRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (designation, description, count, sizes, drawing, order_id) 
		SELECT designation, description, count, sizes, drawing, $1 FROM %s WHERE order_id=$2 ORDER BY id`, OrderPositionTable, OrderPositionTable)

	if _, err := r.db.Exec(query, order.OrderId, order.OldOrderId); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Delete(order *pro_api.DeleteOrderRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", OrdersTable)

	_, err := r.db.Exec(query, order.OrderId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) Save(order *pro_api.SaveOrderRequest) error {
	query := fmt.Sprintf("UPDATE %s SET date=$1 WHERE id=$2", OrdersTable)

	_, err := r.db.Exec(query, time.Now().UnixMilli(), order.OrderId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *OrderRepo) GetPositions(req *pro_api.GetPositionsRequest) (position []models.Position, err error) {
	query := fmt.Sprintf("SELECT id, designation, description, count, sizes, drawing FROM %s WHERE order_id=$1 ORDER BY id", OrderPositionTable)

	if err = r.db.Select(&position, query, req.OrderId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return position, nil
}
