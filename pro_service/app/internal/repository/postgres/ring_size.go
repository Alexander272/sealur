package postgres

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_size_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RingSizeRepo struct {
	db *sqlx.DB
}

func NewRingSizeRepo(db *sqlx.DB) *RingSizeRepo {
	return &RingSizeRepo{
		db: db,
	}
}

type RingSize interface {
	GetAll(context.Context, *ring_size_api.GetRingSize) ([]*ring_size_model.RingSize, error)
	Create(context.Context, *ring_size_api.CreateRingSize) error
	Update(context.Context, *ring_size_api.UpdateRingSize) error
	Delete(context.Context, *ring_size_api.DeleteRingSize) error
}

func (r *RingSizeRepo) GetAll(ctx context.Context, req *ring_size_api.GetRingSize) (sizes []*ring_size_model.RingSize, err error) {
	var data []models.RingSize
	query := fmt.Sprintf(`SELECT id, "outer", "inner" FROM %s ORDER BY "outer", "inner"`, RingSizeTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, rs := range data {
		outer := math.Round(rs.Outer*1000) / 1000
		inner := math.Round(rs.Inner*1000) / 1000
		thickness := math.Ceil((outer - inner) / 2)

		sizes = append(sizes, &ring_size_model.RingSize{
			Id:        rs.Id,
			Outer:     outer,
			Inner:     inner,
			Thickness: thickness,
		})
	}

	return sizes, nil
}

func (r *RingSizeRepo) Create(ctx context.Context, size *ring_size_api.CreateRingSize) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, "outer", "inner", type) VALUES ($1, $2, $3, $4)`, RingSizeTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, size.Outer, size.Inner, size.Type)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingSizeRepo) Update(ctx context.Context, size *ring_size_api.UpdateRingSize) error {
	query := fmt.Sprintf(`UPDATE %s SET "outer"=$1, "inner"=$2, type=$3 WHERE id=$4`, RingSizeTable)

	_, err := r.db.Exec(query, size.Outer, size.Inner, size.Type, size.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingSizeRepo) Delete(ctx context.Context, size *ring_size_api.DeleteRingSize) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RingSizeTable)

	_, err := r.db.Exec(query, size.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
