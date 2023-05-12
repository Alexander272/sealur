package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_construction_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_base_construction_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PutgConstructionBaseRepo struct {
	db *sqlx.DB
}

func NewPutgConstructionBaseRepo(db *sqlx.DB) *PutgConstructionBaseRepo {
	return &PutgConstructionBaseRepo{
		db: db,
	}
}

func (r *PutgConstructionBaseRepo) Get(ctx context.Context, req *putg_base_construction_api.GetPutgBaseConstruction,
) (constructions []*putg_construction_type_model.PutgConstruction, err error) {
	var data []models.PutgConstruction
	query := fmt.Sprintf(`SELECT id, title, code, has_d4, has_d3, has_d2, has_d1, has_rotary_plug, has_inner_ring, has_outer_ring FROM %s ORDER BY code`,
		PutgConstructionTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, c := range data {
		constructions = append(constructions, &putg_construction_type_model.PutgConstruction{
			Id:            c.Id,
			Title:         c.Title,
			Code:          c.Code,
			HasD4:         c.HasD4,
			HasD3:         c.HasD3,
			HasD2:         c.HasD2,
			HasD1:         c.HasD1,
			HasRotaryPlug: c.HasRotaryPlug,
			HasInnerRing:  c.HasInnerRing,
			HasOuterRing:  c.HasOuterRing,
		})
	}

	return constructions, nil
}

func (r *PutgConstructionBaseRepo) Create(ctx context.Context, c *putg_base_construction_api.CreatePutgBaseConstruction) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, title, code, has_d4, has_d3, has_d2, has_d1, has_rotary_plug, has_inner_ring, has_outer_ring)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, PutgConstructionTable,
	)
	id := uuid.New()

	_, err := r.db.Exec(query, id, c.Title, c.Code, c.HasD4, c.HasD3, c.HasD2, c.HasD1, c.HasRotaryPlug, c.HasInnerRing, c.HasOuterRing)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgConstructionBaseRepo) Update(ctx context.Context, c *putg_base_construction_api.UpdatePutgBaseConstruction) error {
	query := fmt.Sprintf(`UPDATE %s SET title=$1, code=$2, has_d4=$3, has_d3=$4, has_d2=$5, has_d1=$6, has_rotary=$7, has_inner_ring=$8,
		has_outer_ring=$9 WHERE id=$10`, PutgConstructionTable,
	)

	_, err := r.db.Exec(query, c.Title, c.Code, c.HasD4, c.HasD3, c.HasD2, c.HasD1, c.HasRotaryPlug, c.HasInnerRing, c.HasOuterRing)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgConstructionBaseRepo) Delete(ctx context.Context, c *putg_base_construction_api.DeletePutgBaseConstruction) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgConstructionTable)

	_, err := r.db.Exec(query, c.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
