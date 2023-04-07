package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_size_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_size_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SnpSizeRepo struct {
	db *sqlx.DB
}

func NewSnpSizeRepo(db *sqlx.DB) *SnpSizeRepo {
	return &SnpSizeRepo{db: db}
}

func (r *SnpSizeRepo) Get(ctx context.Context, req *snp_size_api.GetSnpSize) (sizes []*snp_size_model.SnpSize, err error) {
	var data []models.SnpSize
	query := fmt.Sprintf(`SELECT id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, s2, s3
		FROM %s WHERE snp_type_id=$1 ORDER BY count`, SnpSizeTable)

	if err := r.db.Select(&data, query, req.TypeId); err != nil {
		return nil, fmt.Errorf("failed to execute query")
	}

	for i, ss := range data {
		Pn := []*snp_size_model.Pn{}
		for _, v := range ss.PnMpa {
			Pn = append(Pn, &snp_size_model.Pn{
				Mpa: v,
			})
		}
		for j, v := range ss.PnKg {
			Pn[j].Kg = v
		}

		if i > 0 && ss.Dn == sizes[len(sizes)-1].Dn {
			sizes[len(sizes)-1].Sizes = append(sizes[len(sizes)-1].Sizes, &snp_size_model.Size{
				Pn: Pn,
				D4: ss.D4,
				D3: ss.D3,
				D2: ss.D2,
				D1: ss.D1,
				H:  ss.H,
				S3: ss.S3,
				S2: ss.S2,
			})
		} else {
			sizes = append(sizes, &snp_size_model.SnpSize{
				Id:   ss.Id,
				Dn:   ss.Dn,
				DnMm: ss.DnMm,
				Sizes: []*snp_size_model.Size{{
					Pn: Pn,
					D4: ss.D4,
					D3: ss.D3,
					D2: ss.D2,
					D1: ss.D1,
					H:  ss.H,
					S3: ss.S3,
					S2: ss.S2,
				}},
			})

			if req.HasD2 {
				sizes[len(sizes)-1].D2 = ss.D2
			}
		}
	}

	return sizes, nil
}

func (r *SnpSizeRepo) Create(ctx context.Context, size *snp_size_api.CreateSnpSize) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, snp_type_id, count, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, s2, s3)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, SnpSizeTable)
	id := uuid.New()

	pnMpa := pq.StringArray{}
	pnKg := pq.StringArray{}

	for _, p := range size.Pn {
		pnMpa = append(pnMpa, p.Mpa)
		if p.Kg != "" {
			pnKg = append(pnKg, p.Kg)
		}
	}

	_, err := r.db.Exec(query, id, size.SnpTypeId, size.Count, size.Dn, size.DnMm, pnMpa, pnKg, size.D4, size.D3, size.D2, size.D1,
		pq.Array(size.H), pq.Array(size.S2), pq.Array(size.S3),
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

//TODO написать создание нескольких размеров

func (r *SnpSizeRepo) Update(ctx context.Context, size *snp_size_api.UpdateSnpSize) error {
	query := fmt.Sprintf(`UPDATE %s SET snp_type_id=$1, count=$2, dn=$3, dn_mm=$4, pn_mpa=$5, pn_kg=$6, d4=$7, d3=$8, d2=$9, d1=$10,
		h=$11, s2=$12, s3=$13 WHERE id=$14`, SnpSizeTable,
	)

	pnMpa := pq.StringArray{}
	pnKg := pq.StringArray{}

	for _, p := range size.Pn {
		pnMpa = append(pnMpa, p.Mpa)
		if p.Kg != "" {
			pnKg = append(pnKg, p.Kg)
		}
	}

	_, err := r.db.Exec(query, size.SnpTypeId, size.Count, size.Dn, size.DnMm, pnMpa, pnKg, size.D4, size.D3, size.D2, size.D1,
		pq.Array(size.H), pq.Array(size.S2), pq.Array(size.S3), size.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SnpSizeRepo) Delete(ctx context.Context, size *snp_size_api.DeleteSnpSize) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, SnpSizeTable)

	if _, err := r.db.Exec(query, size.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
