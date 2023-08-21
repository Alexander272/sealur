package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_material_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_material_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RingMaterialRepo struct {
	db *sqlx.DB
}

func NewRingMaterialRepo(db *sqlx.DB) *RingMaterialRepo {
	return &RingMaterialRepo{
		db: db,
	}
}

type RingMaterial interface {
	Get(context.Context, *ring_material_api.GetRingMaterial) ([]*ring_material_model.RingMaterial, error)
	Create(context.Context, *ring_material_api.CreateRingMaterial) error
	Update(context.Context, *ring_material_api.UpdateRingMaterial) error
	Delete(context.Context, *ring_material_api.DeleteRingMaterial) error
}

func (r *RingMaterialRepo) Get(ctx context.Context, req *ring_material_api.GetRingMaterial) (m []*ring_material_model.RingMaterial, err error) {
	var data []models.RingMaterial
	query := fmt.Sprintf(`SELECT id, type, title, description, is_default FROM %s WHERE type=$1 ORDER BY count`, RingMaterialTable)

	if err := r.db.Select(&data, query, req.Type); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, rt := range data {
		m = append(m, &ring_material_model.RingMaterial{
			Id:          rt.Id,
			Title:       rt.Title,
			Type:        rt.Type,
			Description: rt.Description,
			IsDefault:   rt.IsDefault,
		})
	}

	return m, nil
}

func (r *RingMaterialRepo) Create(ctx context.Context, m *ring_material_api.CreateRingMaterial) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, type, title, is_default, count) VALUES ($1, $2, $3, $4, $5)`, RingMaterialTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, m.Type, m.Title, m.IsDefault, m.Count)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingMaterialRepo) Update(ctx context.Context, m *ring_material_api.UpdateRingMaterial) error {
	query := fmt.Sprintf(`UPDATE %s SET type=$1, title=$2, is_default=$3, count=$4 WHERE id=$5`, RingMaterialTable)

	_, err := r.db.Exec(query, m.Type, m.Title, m.IsDefault, m.Count, m.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingMaterialRepo) Delete(ctx context.Context, m *ring_material_api.DeleteRingMaterial) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RingMaterialTable)

	_, err := r.db.Exec(query, m.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
