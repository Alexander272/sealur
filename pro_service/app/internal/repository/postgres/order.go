package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/order_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Get(ctx context.Context, req *order_api.GetOrder) (order *order_model.FullOrder, err error) {
	var data models.ManagerOrder
	query := fmt.Sprintf("SELECT id, date, count_position, number, user_id FROM \"%s\" WHERE id=$1", OrderTable)

	if err := r.db.Get(&data, query, req.Id); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	order = &order_model.FullOrder{
		Id:            data.Id,
		Date:          data.Date,
		CountPosition: data.Count,
		Number:        data.Number,
		UserId:        data.UserId,
	}

	return order, nil
}

func (r *OrderRepo) GetCurrent(ctx context.Context, req *order_api.GetCurrentOrder) (order *order_model.CurrentOrder, err error) {
	var data models.OrderNew
	query := fmt.Sprintf("SELECT id, number FROM \"%s\" WHERE user_id=$1 AND date=''", OrderTable)

	if err := r.db.Get(&data, query, req.UserId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	order = &order_model.CurrentOrder{
		Id:     data.Id,
		Number: data.Number,
	}

	return order, nil
}

func (r *OrderRepo) GetAll(ctx context.Context, req *order_api.GetAllOrders) (orders []*order_model.Order, err error) {
	var data []models.OrderWithPosition
	query := fmt.Sprintf(`SELECT "%s".id, date, count_position, number, %s.id as position_id, title, amount, %s.count as position_count
		FROM "%s" INNER JOIN %s on order_id="%s".id WHERE user_id=$1 AND date != '' ORDER BY number`,
		OrderTable, PositionTable, PositionTable, OrderTable, PositionTable, OrderTable,
	)

	if err := r.db.Select(&data, query, req.UserId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for i, o := range data {
		if i > 0 && o.Id == orders[len(orders)-1].Id {
			orders[len(orders)-1].Positions = append(orders[len(orders)-1].Positions, &position_model.Position{
				Id:     o.PositionId,
				Count:  o.PositionCount,
				Title:  o.Title,
				Amount: o.Amount,
			})
		} else {
			orders = append(orders, &order_model.Order{
				Id:            o.Id,
				Date:          o.Date,
				CountPosition: o.Count,
				Number:        o.Number,
				Positions: []*position_model.Position{{
					Id:     o.PositionId,
					Count:  o.PositionCount,
					Title:  o.Title,
					Amount: o.Amount,
				}},
			})
		}

	}

	return orders, nil
}

func (r *OrderRepo) GetNumber(ctx context.Context, order *order_api.CreateOrder, date string) (int64, error) {
	query := fmt.Sprintf(`UPDATE "%s" SET date=$1, count_position=$2 WHERE id=$3 RETURNING number`, OrderTable)

	row := r.db.QueryRow(query, date, order.Count, order.Id)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed to execute query. error: %w", row.Err())
	}

	var number int64
	if err := row.Scan(&number); err != nil {
		return 0, fmt.Errorf("failed to scan result. error: %w", err)
	}
	return number, nil
}

func (r *OrderRepo) GetOpen(ctx context.Context, managerId string) (orders []*order_model.ManagerOrder, err error) {
	var data []models.ManagerOrder
	// query := fmt.Sprintf(`SELECT "%s".id, user_id, "%s".company, "%s".date, count_position, "number", status FROM "%s"
	// 	INNER JOIN "%s" ON "%s".id=user_id WHERE manager_id=$1 AND status != '%s' AND "%s".date !='' ORDER BY status, "%s".date`,
	// 	OrderTable, UserTable, OrderTable, OrderTable, UserTable, UserTable, order_model.OrderStatus_finish.String(), OrderTable, OrderTable,
	// )
	query := fmt.Sprintf(`SELECT "%s".id, user_id, "%s".company, "%s".date, count_position, "number", status FROM "%s" 
		INNER JOIN "%s" ON "%s".id=user_id WHERE manager_id=$1 AND status != '%s' AND "%s".date !='' ORDER BY "%s".date`,
		OrderTable, UserTable, OrderTable, OrderTable, UserTable, UserTable, order_model.OrderStatus_finish.String(), OrderTable, OrderTable,
	)

	if err := r.db.Select(&data, query, managerId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, o := range data {
		status := order_model.OrderStatus_value[o.Status]

		orders = append(orders, &order_model.ManagerOrder{
			Id:            o.Id,
			Date:          o.Date,
			UserId:        o.UserId,
			Company:       o.Company,
			CountPosition: o.Count,
			Number:        o.Number,
			Status:        order_model.OrderStatus(status),
		})
	}

	return orders, nil
}

func (r *OrderRepo) Create(ctx context.Context, order *order_api.CreateOrder, date string) error {
	query := fmt.Sprintf(`INSERT INTO "%s" (id, user_id, date, count_position) VALUES ($1, $2, $3, $4)`, OrderTable)

	_, err := r.db.Exec(query, order.Id, order.UserId, date, order.Count)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
