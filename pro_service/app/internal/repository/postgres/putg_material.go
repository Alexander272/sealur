package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_material_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PutgMaterialRepo struct {
	db *sqlx.DB
}

func NewPutgMaterialRepo(db *sqlx.DB) *PutgMaterialRepo {
	return &PutgMaterialRepo{
		db: db,
	}
}

func (r *PutgMaterialRepo) Get(ctx context.Context, req *putg_material_api.GetPutgMaterial) (*putg_material_model.PutgMaterials, error) {
	var data []models.SnpMaterial
	query := fmt.Sprintf(`SELECT %s.id, material_id, type, is_default, %s.code, %s.code as base_code, title
		FROM %s INNER JOIN %s ON material_id=%s.id WHERE putg_standard_id=$1 ORDER BY type, count`,
		PutgMaterialTable, PutgMaterialTable, MaterialTable, PutgMaterialTable, MaterialTable, MaterialTable,
	)

	if err := r.db.Select(&data, query, req.StandardId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	var rotaryPlug []*putg_material_model.Material
	var innerRing []*putg_material_model.Material
	var outerRing []*putg_material_model.Material
	var plugDefIndex, innerDefIndex, outerDefIndex int64

	for _, m := range data {
		currentMaterial := &putg_material_model.Material{
			Id:         m.Id,
			MaterialId: m.MaterialId,
			Type:       m.Type,
			IsDefault:  m.IsDefault,
			Code:       m.Code,
			BaseCode:   m.BaseCode,
			Title:      m.Title,
		}

		if m.Type == "rotaryPlug" {
			rotaryPlug = append(rotaryPlug, currentMaterial)
		}
		if m.Type == "innerRing" {
			innerRing = append(innerRing, currentMaterial)
		}
		if m.Type == "outerRing" {
			outerRing = append(outerRing, currentMaterial)
		}
	}

	for i, m := range rotaryPlug {
		if m.IsDefault {
			plugDefIndex = int64(i)
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

	material := &putg_material_model.PutgMaterials{
		RotaryPlug:             rotaryPlug,
		InnerRing:              innerRing,
		OuterRing:              outerRing,
		RotaryPlugDefaultIndex: plugDefIndex,
		InnerRingDefaultIndex:  innerDefIndex,
		OuterRingDefaultIndex:  outerDefIndex,
	}

	return material, nil
}

func (r *PutgMaterialRepo) Create(ctx context.Context, m *putg_material_api.CreatePutgMaterial) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, putg_standard_id, material_id, type, is_default, code) VALUES ($1, $2, $3, $4, $5, $6)`, PutgMaterialTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, m.StandardId, m.MaterialId, m.Type, m.IsDefault, m.Code)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgMaterialRepo) Update(ctx context.Context, m *putg_material_api.UpdatePutgMaterial) error {
	query := fmt.Sprintf(`UPDATE %s SET putg_standard_id=$1, material_id=$2, type=$3, is_default=$4, code=$5 WHERE id=$6`, PutgMaterialTable)

	_, err := r.db.Exec(query, m.StandardId, m.MaterialId, m.Type, m.IsDefault, m.Code, m.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgMaterialRepo) Delete(ctx context.Context, m *putg_material_api.DeletePutgMaterial) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgMaterialTable)

	_, err := r.db.Exec(query, m.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
