package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/jmoiron/sqlx"
)

type PositionRepo struct {
	db *sqlx.DB
}

func NewPositionRepo(db *sqlx.DB) *PositionRepo {
	return &PositionRepo{db: db}
}

func (r *PositionRepo) Get(ctx context.Context) {}

func (r *PositionRepo) CreateSeveral(ctx context.Context, positions []*position_model.Position) error {
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
