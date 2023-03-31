package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/material_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SnpMaterialRepo struct {
	db *sqlx.DB
}

func NewSnpMaterialRepo(db *sqlx.DB) *SnpMaterialRepo {
	return &SnpMaterialRepo{db: db}
}

func (r *SnpMaterialRepo) Get(ctx context.Context, req *snp_material_api.GetSnpMaterial) (materials []*snp_material_model.SnpMaterial, err error) {
	var data []models.SNPMaterial
	query := fmt.Sprintf(`SELECT %s.id, title, code, short_en, short_rus, %s.id as material_id, default_id, type
		FROM %s INNER JOIN %s ON array[%s.id]<@material_id WHERE standard_id=$1 ORDER BY type, count`,
		SnpMaterialTable, MaterialTable, SnpMaterialTable, MaterialTable, MaterialTable,
	)

	if err := r.db.Select(&data, query, req.StandardId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for i, s := range data {
		if i > 0 && s.Id == materials[len(materials)-1].Id {
			materials[len(materials)-1].Materials = append(materials[len(materials)-1].Materials, &material_model.Material{
				Id:       s.MaterialId,
				Title:    s.Title,
				Code:     s.Code,
				ShortEn:  s.ShortEn,
				ShortRus: s.ShortRus,
			})
		} else {
			materials = append(materials, &snp_material_model.SnpMaterial{
				Id:   s.Id,
				Type: s.Type,
				Materials: []*material_model.Material{{
					Id:       s.MaterialId,
					Title:    s.Title,
					Code:     s.Code,
					ShortEn:  s.ShortEn,
					ShortRus: s.ShortRus,
				}},
			})
		}
		if s.Default == s.MaterialId {
			materials[len(materials)-1].Default = &material_model.Material{
				Id:       s.MaterialId,
				Title:    s.Title,
				Code:     s.Code,
				ShortEn:  s.ShortEn,
				ShortRus: s.ShortRus,
			}
		}
	}

	return materials, nil
}

func (r *SnpMaterialRepo) Create(ctx context.Context, material *snp_material_api.CreateSnpMaterial) error {
	query := fmt.Sprintf("INSERT INTO %s (id, material_id, default_id, type, standard_id) VALUES ($1, $2, $3, $4, $5)", SnpMaterialTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, pq.Array(material.MaterialId), material.Default, material.Type, material.StandardId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SnpMaterialRepo) Update(ctx context.Context, material *snp_material_api.UpdateSnpMaterial) error {
	query := fmt.Sprintf("UPDATE %s	SET material_id=$1, default_id=$2, type=$3, standard_id=$4 WHERE id=$5", SnpMaterialTable)

	_, err := r.db.Exec(query, pq.Array(material.MaterialId), material.Default, material.Type, material.StandardId, material.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SnpMaterialRepo) Delete(ctx context.Context, material *snp_material_api.DeleteSnpMaterial) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", SnpMaterialTable)

	if _, err := r.db.Exec(query, material.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
