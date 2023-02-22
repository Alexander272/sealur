package repository

import (
	"context"
	"fmt"

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

	_, err := r.db.Exec(query, pq.Array(material.MaterialId), material.Default, material.Type, material.StandardId)
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
