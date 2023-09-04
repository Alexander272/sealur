package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_construction_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_construction_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RingConstructionRepo struct {
	db *sqlx.DB
}

func NewRingConstructionRepo(db *sqlx.DB) *RingConstructionRepo {
	return &RingConstructionRepo{
		db: db,
	}
}

type RingConstruction interface {
	GetAll(context.Context, *ring_construction_api.GetRingConstructions) (*ring_construction_model.RingConstructionMap, error)
	Create(context.Context, *ring_construction_api.CreateRingConstruction) error
	Update(context.Context, *ring_construction_api.UpdateRingConstruction) error
	Delete(context.Context, *ring_construction_api.DeleteRingConstruction) error
}

func (r *RingConstructionRepo) GetAll(ctx context.Context, req *ring_construction_api.GetRingConstructions,
) (*ring_construction_model.RingConstructionMap, error) {
	var data []models.RingConstruction
	query := fmt.Sprintf(`SELECT id, type_id, code, base_code, title, description, image, without_rotary_plug FROM %s ORDER BY count`, RingConstructionTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	constructions := make(map[string]*ring_construction_model.Construction, 0)

	for _, rd := range data {
		c := &ring_construction_model.RingConstruction{
			Id:                rd.Id,
			Code:              rd.Code,
			Title:             rd.Title,
			Description:       rd.Description,
			Image:             rd.Image,
			WithoutRotaryPlug: rd.WithoutRotaryPlug,
			BaseCode:          rd.BaseCode,
		}

		if constructions[rd.TypeId] == nil {
			constructions[rd.TypeId] = &ring_construction_model.Construction{Constructions: []*ring_construction_model.RingConstruction{c}}
		} else {
			constructions[rd.TypeId].Constructions = append(constructions[rd.TypeId].Constructions, c)
		}
	}

	return &ring_construction_model.RingConstructionMap{Constructions: constructions}, nil
}

func (r *RingConstructionRepo) Create(ctx context.Context, c *ring_construction_api.CreateRingConstruction) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, type_id, code, title, description, image, count, without_rotary_plug) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		RingConstructionTable,
	)
	id := uuid.New()

	_, err := r.db.Exec(query, id, c.TypeId, c.Code, c.Title, c.Description, c.Image, c.Count, c.WithoutRotaryPlug)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingConstructionRepo) Update(ctx context.Context, c *ring_construction_api.UpdateRingConstruction) error {
	query := fmt.Sprintf(`UPDATE %s SET type_id=$1, code=$2, title=$3, description=$4, image=$5, count=$6, without_rotary_plug=$7 WHERE id=$8`,
		RingConstructionTable,
	)

	_, err := r.db.Exec(query, c.TypeId, c.Code, c.Title, c.Description, c.Image, c.Count, c.WithoutRotaryPlug, c.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingConstructionRepo) Delete(ctx context.Context, c *ring_construction_api.DeleteRingConstruction) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RingConstructionTable)

	_, err := r.db.Exec(query, c.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
