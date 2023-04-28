package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_construction_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_construction_api"
	"github.com/jmoiron/sqlx"
)

type PutgConstructionRepo struct {
	db *sqlx.DB
}

func NewPutgConstructionRepo(db *sqlx.DB) *PutgConstructionRepo {
	return &PutgConstructionRepo{
		db: db,
	}
}

func (r *PutgConstructionRepo) Get(ctx context.Context, req *putg_construction_api.GetPutgConstruction,
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

// TODO дописать оставшиеся функции
