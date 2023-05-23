package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_size_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (r *PutgSizeRepo) GetNew(ctx context.Context, req *putg_size_api.GetPutgSize_New) (sizes []*putg_size_model.PutgSize, err error) {
	var data []models.PutgSize
	query := fmt.Sprintf(`SELECT id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h FROM %s
		WHERE putg_flange_type_id=$1 AND base_construction_id=$2 ORDER BY count`, PutgSizeTableTest,
	)

	if err := r.db.Select(&data, query, req.FlangeTypeId, req.BaseConstructionId); err != nil {
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

func (r *PutgSizeRepo) Create(ctx context.Context, size *putg_size_api.CreatePutgSize) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, putg_standard_id, construction_id, count, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`, PutgSizeTable)
	id := uuid.New()

	pnMpa := pq.StringArray{}
	pnKg := pq.StringArray{}

	for _, p := range size.Pn {
		pnMpa = append(pnMpa, p.Mpa)
		if p.Kg != "" {
			pnKg = append(pnKg, p.Kg)
		}
	}

	_, err := r.db.Exec(query, id, size.PutgStandardId, size.ConstructionId, size.Count, size.Dn, size.DnMm, pnMpa, pnKg, size.D4,
		size.D3, size.D2, size.D1, pq.Array(size.H),
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgSizeRepo) CreateSeveral(ctx context.Context, sizes *putg_size_api.CreateSeveralPutgSize) error {
	query := fmt.Sprintf("INSERT INTO %s (id, putg_standard_id, construction_id, count, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h) VALUES ", PutgSizeTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(sizes.Sizes))

	c := 13
	for i, s := range sizes.Sizes {
		id := uuid.New()
		pnMpa := pq.StringArray{}
		pnKg := pq.StringArray{}

		for _, p := range s.Pn {
			pnMpa = append(pnMpa, p.Mpa)
			if p.Kg != "" {
				pnKg = append(pnKg, p.Kg)
			}
		}

		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*c+1, i*c+2, i*c+3, i*c+4, i*c+5, i*c+6, i*c+7, i*c+8, i*c+9, i*c+10, i*c+11, i*c+12, i*c+13,
		))
		args = append(args, id, s.PutgStandardId, s.ConstructionId, s.Count, s.Dn, s.DnMm, pnMpa, pnKg, s.D4, s.D3, s.D2, s.D1, pq.Array(s.H))
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PutgSizeRepo) Update(ctx context.Context, size *putg_size_api.UpdatePutgSize) error {
	query := fmt.Sprintf(`UPDATE %s SET putg_standard_id=$1, count=$2, dn=$3, dn_mm=$4, pn_mpa=$5, pn_kg=$6, d4=$7, d3=$8, d2=$9, d1=$10,
		h=$11, construction_id=$12 WHERE id=$13`, PutgSizeTable,
	)

	pnMpa := pq.StringArray{}
	pnKg := pq.StringArray{}

	for _, p := range size.Pn {
		pnMpa = append(pnMpa, p.Mpa)
		if p.Kg != "" {
			pnKg = append(pnKg, p.Kg)
		}
	}

	_, err := r.db.Exec(query, size.PutgStandardId, size.Count, size.Dn, size.DnMm, pnMpa, pnKg, size.D4, size.D3, size.D2, size.D1,
		pq.Array(size.H), size.ConstructionId, size.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PutgSizeRepo) Delete(ctx context.Context, size *putg_size_api.DeletePutgSize) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgSizeTable)

	if _, err := r.db.Exec(query, size.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
