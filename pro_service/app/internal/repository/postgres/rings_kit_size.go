package postgres

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/rings_kit_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_size_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RingsKitSizeRepo struct {
	db *sqlx.DB
}

func NewRingsKitSizeRepo(db *sqlx.DB) *RingsKitSizeRepo {
	return &RingsKitSizeRepo{db: db}
}

type RingsKitSize interface {
	Get(context.Context, *rings_kit_size_api.GetRingsKitSize) ([]*rings_kit_size_model.RingsKitSize, error)
	Create(context.Context, *rings_kit_size_api.CreateRingsKitSize) error
	Update(context.Context, *rings_kit_size_api.UpdateRingsKitSize) error
	Delete(context.Context, *rings_kit_size_api.DeleteRingsKitSize) error
}

func (r *RingsKitSizeRepo) Get(ctx context.Context, req *rings_kit_size_api.GetRingsKitSize) (sizes []*rings_kit_size_model.RingsKitSize, err error) {
	var data []models.RingsKitSize

	query := fmt.Sprintf(`SELECT id, "outer", "inner", thickness FROM %s WHERE construction_id=$1 ORDER BY "outer", "inner"`, RingsKitSizeTable)

	if err := r.db.Select(&data, query, req.ConstructionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, rs := range data {
		sizes = append(sizes, &rings_kit_size_model.RingsKitSize{
			Id:        rs.Id,
			Outer:     math.Round(rs.Outer*1000) / 1000,
			Inner:     math.Round(rs.Inner*1000) / 1000,
			Thickness: math.Round(rs.Thickness*1000) / 1000,
		})
	}

	return sizes, nil
}

func (r *RingsKitSizeRepo) Create(ctx context.Context, size *rings_kit_size_api.CreateRingsKitSize) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, construction_id, "outer", "inner", thickness) VALUES ($1, $2, $3, $4, $5)`, RingsKitSizeTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, size.ConstructionId, size.Outer, size.Inner, size.Thickness)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingsKitSizeRepo) Update(ctx context.Context, size *rings_kit_size_api.UpdateRingsKitSize) error {
	query := fmt.Sprintf(`UPDATE %s SET construction_id=$1, "outer"=$2, "inner"=$3, thickness=$4 WHERE id=$5`, RingsKitSizeTable)

	_, err := r.db.Exec(query, size.ConstructionId, size.Outer, size.Inner, size.Thickness)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingsKitSizeRepo) Delete(ctx context.Context, size *rings_kit_size_api.DeleteRingsKitSize) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RingsKitSizeTable)

	_, err := r.db.Exec(query, size.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
