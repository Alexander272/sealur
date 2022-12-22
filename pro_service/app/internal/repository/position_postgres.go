package repository

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type PositionRepo struct {
	db *sqlx.DB
}

func NewPositionRepo(db *sqlx.DB) *PositionRepo {
	return &PositionRepo{db: db}
}

func (r *PositionRepo) Get(req *pro_api.GetPositionsRequest) (position []models.Position, err error) {
	query := fmt.Sprintf("SELECT id, designation, description, count, sizes, drawing FROM %s WHERE order_id=$1 ORDER BY id", OrderPositionTable)

	if err = r.db.Select(&position, query, req.OrderId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return position, nil
}

func (r *PositionRepo) GetCur(req *pro_api.GetCurPositionsRequest) (position []models.Position, err error) {
	query := fmt.Sprintf(`SELECT id, designation, description, count, sizes, drawing, order_id FROM %s WHERE order_id=(
		SELECT id FROM %s WHERE user_id=$1 AND date IS NULL
	) ORDER BY id`, OrderPositionTable, OrdersTable)

	if err = r.db.Select(&position, query, req.UserId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return position, nil
}

func (r *PositionRepo) Add(position *pro_api.AddPositionRequest) (id string, err error) {
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

	return fmt.Sprintf("%d", idInt), tx.Commit()
}

func (r *PositionRepo) Update(position *pro_api.UpdatePositionRequest) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET count=$1 WHERE id=$2", OrderPositionTable)
	_, err := r.db.Exec(updateQuery, position.Count, position.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PositionRepo) Remove(position *pro_api.RemovePositionRequest) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction. error: %w", err)
	}

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING drawing", OrderPositionTable)
	row := tx.QueryRow(deleteQuery, position.Id)

	var drawing string
	if err := row.Scan(&drawing); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET count_position=count_position-1 WHERE id=$1", OrdersTable)
	_, err = tx.Exec(updateQuery, position.OrderId)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	return drawing, tx.Commit()
}
