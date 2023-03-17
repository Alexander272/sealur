package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PositionRepo struct {
	db *sqlx.DB
}

func NewPositionRepo(db *sqlx.DB) *PositionRepo {
	return &PositionRepo{db: db}
}

func (r *PositionRepo) Get(ctx context.Context, orderId string) (positions []*position_model.OrderPosition, err error) {
	var data []models.PositionNew
	query := fmt.Sprintf(`SELECT id, title, amount, type, count FROM %s WHERE order_id=$1 ORDER BY count`, PositionTable)

	if err := r.db.Select(&data, query, orderId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, p := range data {
		positionType := position_model.PositionType_value[p.Type]
		positions = append(positions, &position_model.OrderPosition{
			Id:     p.Id,
			Title:  p.Title,
			Count:  p.Count,
			Amount: p.Amount,
			Type:   position_model.PositionType(positionType),
		})
	}

	return positions, nil
}

func (r *PositionRepo) Create(ctx context.Context, position *position_model.FullPosition) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s (id, order_id, title, amount, type, count) VALUES ($1, $2, $3, $4, $5, $6)", PositionTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, position.OrderId, position.Title, position.Amount, position.Type.String(), position.Count)
	if err != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}
	return id.String(), nil
}

func (r *PositionRepo) CreateSeveral(ctx context.Context, positions []*position_model.FullPosition) error {
	query := fmt.Sprintf("INSERT INTO %s (id, order_id, title, amount, type, count) VALUES ", PositionTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(positions))

	c := 6
	for i, s := range positions {
		// id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5, i*c+6))
		args = append(args, s.Id, s.OrderId, s.Title, s.Amount, s.Type.String(), s.Count)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PositionRepo) Update(ctx context.Context, position *position_model.FullPosition) error {
	query := fmt.Sprintf(`UPDATE %s SET title=$1, amount=$2 WHERE id=$3`, PositionTable)

	_, err := r.db.Exec(query, position.Title, position.Amount, position.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PositionRepo) Delete(ctx context.Context, positionId string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", PositionTable)

	_, err := r.db.Exec(query, positionId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}