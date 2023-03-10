package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Get(ctx context.Context) {}

func (r *OrderRepo) GetAll(ctx context.Context) {}

func (r *OrderRepo) GetNumber(ctx context.Context, orderId, date string) (int64, error) {
	query := fmt.Sprintf(`UPDATE %s	SET date=$1	WHERE id=$2 RETURNING number`, OrderTable)

	row, err := r.db.Query(query, orderId, date)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query. error: %w", err)
	}

	var number int64
	if err := row.Scan(&number); err != nil {
		return 0, fmt.Errorf("failed to scan result. error: %w", err)
	}
	return number, nil
}

func (r *OrderRepo) Create(ctx context.Context, order order_api.CreateOrder, date string) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, user_id, date, count_position) VALUES ($1, $2, $3, $4)`, OrderTable)

	_, err := r.db.Exec(query, order.Id, order.UserId, date, order.Count)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
