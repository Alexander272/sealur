package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_material_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SnpMaterialRepo struct {
	db *sqlx.DB
}

func NewSnpMaterialRepo(db *sqlx.DB) *SnpMaterialRepo {
	return &SnpMaterialRepo{db: db}
}

func (r *SnpMaterialRepo) Get(ctx context.Context, req *snp_material_api.GetSnpMaterial) (*snp_material_model.SnpMaterials, error) {
	var data []models.SnpMaterial
	query := fmt.Sprintf(`SELECT %s.id, material_id, type, is_default, %s.code, is_standard, %s.code as base_code, title
		FROM %s INNER JOIN %s ON material_id=%s.id WHERE standard_id=$1 ORDER BY type, count`,
		SnpMaterialTableNew, SnpMaterialTableNew, MaterialTable, SnpMaterialTableNew, MaterialTable, MaterialTable,
	)

	if err := r.db.Select(&data, query, req.StandardId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	var frame []*snp_material_model.Material
	var innerRing []*snp_material_model.Material
	var outerRing []*snp_material_model.Material
	var frameDefIndex, innerDefIndex, outerDefIndex int64

	for _, m := range data {
		currentMaterial := &snp_material_model.Material{
			Id:         m.Id,
			MaterialId: m.MaterialId,
			Type:       m.Type,
			IsDefault:  m.IsDefault,
			Code:       m.Code,
			IsStandard: m.IsStandard,
			BaseCode:   m.BaseCode,
			Title:      m.Title,
		}

		if m.Type == "fr" {
			frame = append(frame, currentMaterial)
		}
		if m.Type == "ir" {
			innerRing = append(innerRing, currentMaterial)
		}
		if m.Type == "or" {
			outerRing = append(outerRing, currentMaterial)
		}
	}

	for i, m := range frame {
		if m.IsDefault {
			frameDefIndex = int64(i)
			break
		}
	}
	for i, m := range innerRing {
		if m.IsDefault {
			innerDefIndex = int64(i)
			break
		}
	}
	for i, m := range outerRing {
		if m.IsDefault {
			outerDefIndex = int64(i)
			break
		}
	}

	material := &snp_material_model.SnpMaterials{
		Frame:                 frame,
		InnerRing:             innerRing,
		OuterRing:             outerRing,
		FrameDefaultIndex:     frameDefIndex,
		InnerRingDefaultIndex: innerDefIndex,
		OuterRingDefaultIndex: outerDefIndex,
	}

	return material, nil
}

func (r *SnpMaterialRepo) Create(ctx context.Context, material *snp_material_api.CreateSnpMaterial) error {
	query := fmt.Sprintf("INSERT INTO %s (id, standard_id, material_id, type, is_default, code, is_standard) VALUES ($1, $2, $3, $4, $5, $6, $7)", SnpMaterialTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, material.StandardId, material.MaterialId, material.Type, material.IsDefault, material.Code, material.IsStandard)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SnpMaterialRepo) Update(ctx context.Context, material *snp_material_api.UpdateSnpMaterial) error {
	query := fmt.Sprintf("UPDATE %s	SET standard_id=$1, material_id=$2, type=$3, is_default=$4, code=$5, is_standard=$6 WHERE id=$7", SnpMaterialTable)

	_, err := r.db.Exec(query, material.StandardId, material.MaterialId, material.Type, material.IsDefault, material.Code, material.IsStandard, material.Id)
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
