package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_size_api"
	"github.com/jmoiron/sqlx"
)

type PutgSizeRepo struct {
	db *sqlx.DB
}

func NewPutgSizeRepo(db *sqlx.DB) *PutgSizeRepo {
	return &PutgSizeRepo{
		db: db,
	}
}

func (r *PutgSizeRepo) Get(ctx context.Context, req *putg_size_api.GetPutgSize) (sizes []*putg_size_model.PutgSize, err error) {
	var data []models.PutgSize
	query := fmt.Sprintf(`SELECT id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h FROM %s 
		WHERE putg_standard_id=$1 AND construction_id=$2 ORDER BY count`, PutgSizeTable,
	)

	if err := r.db.Select(&data, query, req.PutgStandardId, req.ConstructionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for i, ps := range data {
		Pn := []*putg_size_model.Pn{}
		for _, v := range ps.PnMpa {
			Pn = append(Pn, &putg_size_model.Pn{
				Mpa: v,
			})
		}
		for j, v := range ps.PnKg {
			Pn[j].Kg = v
		}

		if i > 0 && ps.Dn == sizes[len(sizes)-1].Dn {
			sizes[len(sizes)-1].Sizes = append(sizes[len(sizes)-1].Sizes, &putg_size_model.Size{
				Pn: Pn,
				D4: ps.D4,
				D3: ps.D3,
				D2: ps.D2,
				D1: ps.D1,
				H:  ps.H,
			})
		} else {
			sizes = append(sizes, &putg_size_model.PutgSize{
				Id:   ps.Id,
				Dn:   ps.Dn,
				DnMm: ps.DnMm,
				Sizes: []*putg_size_model.Size{{
					Pn: Pn,
					D4: ps.D4,
					D3: ps.D3,
					D2: ps.D2,
					D1: ps.D1,
					H:  ps.H,
				}},
			})
		}
	}

	return sizes, nil
}

// TODO дописать оставшиеся функции
