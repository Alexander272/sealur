package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/material_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/material_model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MaterialRepo struct {
	db *sqlx.DB
}

func NewMaterialRepo(db *sqlx.DB) *MaterialRepo {
	return &MaterialRepo{db: db}
}

func (r *MaterialRepo) GetAll(ctx context.Context, mat *material_api.GetAllMaterials) (materials []*material_model.Material, err error) {
	var data []models.Material
	query := fmt.Sprintf("SELECT id, title, code, short_en, short_rus FROM %s", MaterialTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, m := range data {
		materials = append(materials, &material_model.Material{
			Id:       m.Id,
			Title:    m.Title,
			Code:     m.Code,
			ShortEn:  m.ShortEn,
			ShortRus: m.ShortRus,
		})
	}

	return materials, nil
}

func (r *MaterialRepo) Create(ctx context.Context, material *material_api.CreateMaterial) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code, short_en, short_rus) VALUES ($1, $2, $3, $4, $5)", MaterialTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, material.Title, material.Code, material.ShortEn, material.ShortRus)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialRepo) CreateSeveral(ctx context.Context, materials *material_api.CreateSeveralMaterial) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code, short_en, short_rus) VALUES ", MaterialTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(materials.Materials))

	c := 5
	for i, m := range materials.Materials {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5))
		args = append(args, id, m.Title, m.Code, m.ShortEn, m.ShortRus)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialRepo) Update(ctx context.Context, material *material_api.UpdateMaterial) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, code=$2, short_en=$3, short_rus=$4 WHERE id=$5", MaterialTable)

	_, err := r.db.Exec(query, material.Title, material.Code, material.ShortEn, material.ShortRus, material.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialRepo) Delete(ctx context.Context, material *material_api.DeleteMaterial) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", MaterialTable)

	if _, err := r.db.Exec(query, material.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
