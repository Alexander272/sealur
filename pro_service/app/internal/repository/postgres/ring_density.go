package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_density_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_density_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RingDensityRepo struct {
	db *sqlx.DB
}

func NewRingDensityRepo(db *sqlx.DB) *RingDensityRepo {
	return &RingDensityRepo{
		db: db,
	}
}

type RingDensity interface {
	GetAll(context.Context, *ring_density_api.GetRingDensity) (*ring_density_model.RingDensityMap, error)
	Create(context.Context, *ring_density_api.CreateRingDensity) error
	Update(context.Context, *ring_density_api.UpdateRingDensity) error
	Delete(context.Context, *ring_density_api.DeleteRingDensity) error
}

func (r *RingDensityRepo) GetAll(ctx context.Context, req *ring_density_api.GetRingDensity) (*ring_density_model.RingDensityMap, error) {
	var data []models.RingDensity
	query := fmt.Sprintf(`SELECT id, type_id, code, title, description, has_rotary_plug FROM %s ORDER BY type_id, title`, RingDensityTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	density := make(map[string]*ring_density_model.Density, 0)

	for _, rd := range data {
		d := &ring_density_model.RingDensity{
			Id:            rd.Id,
			Code:          rd.Code,
			Title:         rd.Title,
			Description:   rd.Description,
			HasRotaryPlug: rd.HasRotaryPlug,
		}

		if density[rd.TypeId] == nil {
			density[rd.TypeId] = &ring_density_model.Density{Density: []*ring_density_model.RingDensity{d}}
		} else {
			density[rd.TypeId].Density = append(density[rd.TypeId].Density, d)
		}
	}

	return &ring_density_model.RingDensityMap{Density: density}, nil
}

func (r *RingDensityRepo) Create(ctx context.Context, density *ring_density_api.CreateRingDensity) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, type_id, code, title, description, has_rotary_plug)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		RingDensityTable,
	)
	id := uuid.New()

	_, err := r.db.Exec(query, id, density.TypeId, density.Code, density.Title, density.Description, density.HasRotaryPlug)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingDensityRepo) Update(ctx context.Context, density *ring_density_api.UpdateRingDensity) error {
	query := fmt.Sprintf(`UPDATE %s SET type_id=$1, code=$2, title=$3, description=$4, has_rotary_plug=$5 WHERE id=$6`, RingDensityTable)

	_, err := r.db.Exec(query, density.TypeId, density.Code, density.Title, density.Description, density.HasRotaryPlug, density.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingDensityRepo) Delete(ctx context.Context, density *ring_density_api.DeleteRingDensity) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RingDensityTable)

	_, err := r.db.Exec(query, density.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
