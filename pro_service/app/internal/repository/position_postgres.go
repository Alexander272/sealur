package repository

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type PositionRepo struct {
	db *sqlx.DB
}

func NewPositionRepo(db *sqlx.DB) *PositionRepo {
	return &PositionRepo{db: db}
}

func (r *PositionRepo) Get(req *proto.GetPositionsRequest) (position []models.Position, err error) {
	query := fmt.Sprintf("SELECT id, designation, description, count, sizes, drawing FROM %s WHERE order_id=$1 ORDER BY id", OrderPositionTable)

	if err = r.db.Select(&position, query, req.OrderId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return position, nil
}

func (r *PositionRepo) GetCur(req *proto.GetCurPositionsRequest) (position []models.Position, err error) {
	query := fmt.Sprintf(`SELECT id, designation, description, count, sizes, drawing, order_id FROM %s WHERE order_id=(
		SELECT id FROM %s WHERE user_id=$1 AND date IS NULL
	) ORDER BY id`, OrderPositionTable, OrdersTable)

	if err = r.db.Select(&position, query, req.UserId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return position, nil
}

func (r *PositionRepo) Add(position *proto.AddPositionRequest) (id string, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return id, fmt.Errorf("failed to start transaction. error: %w", err)
	}

	createQuery := fmt.Sprintf(`INSERT INTO %s (designation, description, count, sizes, drawing, order_id) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, OrderPositionTable)
	row := tx.QueryRow(createQuery, position.Designation, position.Description, position.Count, position.Sizes, position.Drawing, position.OrderId)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		tx.Rollback()
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET count_position=count_position+1 WHERE id=$1", OrdersTable)
	_, err = tx.Exec(updateQuery, position.OrderId)
	if err != nil {
		tx.Rollback()
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, tx.Commit()
}

func (r *PositionRepo) Update(position *proto.UpdatePositionRequest) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET count=$1 WHERE id=$2", OrderPositionTable)
	_, err := r.db.Exec(updateQuery, position.Count, position.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PositionRepo) Remove(position *proto.RemovePositionRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction. error: %w", err)
	}

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id=$1", OrderPositionTable)
	_, err = tx.Exec(deleteQuery, position.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET count_position=count_position-1 WHERE id=$1", OrdersTable)
	_, err = tx.Exec(updateQuery, position.OrderId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return tx.Commit()
}
