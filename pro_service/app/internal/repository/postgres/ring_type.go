package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_type_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RingTypeRepo struct {
	db *sqlx.DB
}

func NewRingTypeRepo(db *sqlx.DB) *RingTypeRepo {
	return &RingTypeRepo{
		db: db,
	}
}

type RingType interface {
	GetAll(context.Context, *ring_type_api.GetRingTypes) ([]*ring_type_model.RingType, error)
	Create(context.Context, *ring_type_api.CreateRingType) error
	Update(context.Context, *ring_type_api.UpdateRingType) error
	Delete(context.Context, *ring_type_api.DeleteRingType) error
}

func (r *RingTypeRepo) GetAll(ctx context.Context, req *ring_type_api.GetRingTypes) (ringTypes []*ring_type_model.RingType, err error) {
	var data []models.RingType
	query := fmt.Sprintf(`SELECT id, code, title, description, has_rotary_plug, has_density, has_thickness, material_type, image FROM %s
		ORDER BY count`,
		RingTypeTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, rt := range data {
		ringTypes = append(ringTypes, &ring_type_model.RingType{
			Id:            rt.Id,
			Code:          rt.Code,
			Title:         rt.Title,
			Description:   rt.Description,
			HasRotaryPlug: rt.HasRotaryPlug,
			HasDensity:    rt.HasDensity,
			HasThickness:  rt.HasThickness,
			MaterialType:  rt.MaterialType,
			// Image:         rt.Image,
		})
	}

	return ringTypes, nil
}

func (r *RingTypeRepo) Create(ctx context.Context, ring *ring_type_api.CreateRingType) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, code, title, description, has_rotary_plug, has_density, has_thickness, material_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		RingTypeTable,
	)
	id := uuid.New()

	_, err := r.db.Exec(query, id, ring.Code, ring.Title, ring.Description, ring.HasRotaryPlug, ring.HasDensity, ring.HasThickness, ring.MaterialType)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingTypeRepo) Update(ctx context.Context, ring *ring_type_api.UpdateRingType) error {
	query := fmt.Sprintf(`UPDATE %s SET code=$1, title=$2, description=$3, has_rotary_plug=$4, has_density=$5, has_thickness=$6, material_type=$7
		WHERE id=$8`,
		RingTypeTable,
	)

	_, err := r.db.Exec(query, ring.Code, ring.Title, ring.Description, ring.HasRotaryPlug, ring.HasDensity, ring.HasThickness, ring.MaterialType, ring.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingTypeRepo) Delete(ctx context.Context, ring *ring_type_api.DeleteRingType) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RingTypeTable)

	_, err := r.db.Exec(query, ring.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
