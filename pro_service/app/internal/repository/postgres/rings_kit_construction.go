package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/rings_kit_construction_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_construction_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RingsKitConstructionRepo struct {
	db *sqlx.DB
}

func NewRingsKitConstructionRepo(db *sqlx.DB) *RingsKitConstructionRepo {
	return &RingsKitConstructionRepo{
		db: db,
	}
}

type RingsKitConstruction interface {
	GetAll(context.Context, *rings_kit_construction_api.GetRingsKitConstructions) (*rings_kit_construction_model.RingsKitConstructionMap, error)
	Create(context.Context, *rings_kit_construction_api.CreateRingsKitConstruction) error
	Update(context.Context, *rings_kit_construction_api.UpdateRingsKitConstruction) error
	Delete(context.Context, *rings_kit_construction_api.DeleteRingsKitConstruction) error
}

func (r *RingsKitConstructionRepo) GetAll(ctx context.Context, req *rings_kit_construction_api.GetRingsKitConstructions,
) (*rings_kit_construction_model.RingsKitConstructionMap, error) {
	var data []models.RingsKitConstruction
	query := fmt.Sprintf(`SELECT id, type_id, code, title, image, same_rings, material_types, has_thickness, default_count, default_materials 
		FROM %s ORDER BY count`,
		RingsKitConstructionTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	constructions := make(map[string]*rings_kit_construction_model.KitConstruction, 0)

	for _, rd := range data {
		c := &rings_kit_construction_model.RingsKitConstruction{
			Id:               rd.Id,
			Code:             rd.Code,
			Title:            rd.Title,
			Image:            rd.Image,
			SameRings:        rd.SameRings,
			MaterialTypes:    rd.MaterialTypes,
			HasThickness:     rd.HasThickness,
			DefaultCount:     rd.DefaultCount,
			DefaultMaterials: rd.DefaultMaterials,
		}

		if constructions[rd.TypeId] == nil {
			constructions[rd.TypeId] = &rings_kit_construction_model.KitConstruction{Constructions: []*rings_kit_construction_model.RingsKitConstruction{c}}
		} else {
			constructions[rd.TypeId].Constructions = append(constructions[rd.TypeId].Constructions, c)
		}
	}

	return &rings_kit_construction_model.RingsKitConstructionMap{Constructions: constructions}, nil
}

func (r *RingsKitConstructionRepo) Create(ctx context.Context, c *rings_kit_construction_api.CreateRingsKitConstruction) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, type_id, code, title, image, count, same_rings, material_types, has_thickness, default_count, default_materials)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		RingsKitConstructionTable,
	)
	id := uuid.New()

	_, err := r.db.Exec(query, id, c.TypeId, c.Code, c.Title, c.Image, c.Count, c.SameRings, c.MaterialTypes, c.HasThickness, c.DefaultCount, c.DefaultMaterials)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingsKitConstructionRepo) Update(ctx context.Context, c *rings_kit_construction_api.UpdateRingsKitConstruction) error {
	query := fmt.Sprintf(`UPDATE %s SET type_id=$1, code=$2, title=$3, image=$4, count=$5, same_rings=$6, material_types=$7, has_thickness=$8, 
		default_count=$9, default_materials=$10 WHERE id=$11`,
		RingsKitConstructionTable,
	)

	_, err := r.db.Exec(query, c.TypeId, c.Code, c.Title, c.Image, c.Count, c.SameRings, c.MaterialTypes, c.HasThickness, c.DefaultCount, c.DefaultMaterials,
		c.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingsKitConstructionRepo) Delete(ctx context.Context, c *rings_kit_construction_api.DeleteRingsKitConstruction) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RingsKitConstructionTable)

	_, err := r.db.Exec(query, c.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
