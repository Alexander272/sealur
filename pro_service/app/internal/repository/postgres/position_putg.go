package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PositionPutgRepo struct {
	db *sqlx.DB
}

func NewPositionPutgRepo(db *sqlx.DB) *PositionPutgRepo {
	return &PositionPutgRepo{
		db: db,
	}
}

func (r *PositionPutgRepo) Create(ctx context.Context, position *position_model.FullPosition) error {
	mainQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, putg_standard_id, flange_type_id, configuration_id)
		VALUES ($1, $2, $3, $4, $5)`, PositionMainPutgTable,
	)
	sizeQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, dn, dn_mm, pn_mpa, pn_kg, d4, d3, d2, d1, h, another, use_dimensions)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`, PositionSizePutgTable,
	)
	materialQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, filler_id, filler_code, type_id, construction_id, construction_code, 
		rotary_plug_id, rotary_plug_code, inner_ring_id, inner_ring_code, outer_ring_id, outer_ring_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`, PositionMaterialPutgTable,
	)
	designQuery := fmt.Sprintf(`INSERT INTO %s (id, position_id, has_jumper, jumper_code, jumper_width, has_hole, has_removable, 
		has_mounting, mounting_code, drawing) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, PositionDesignPutgTable,
	)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction. error: %w", err)
	}

	id := uuid.New()
	main := position.PutgData.Main
	size := position.PutgData.Size
	material := position.PutgData.Material
	design := position.PutgData.Design

	_, err = tx.Exec(mainQuery, id, position.Id, main.PutgStandardId, main.FlangeTypeId, main.ConfigurationId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query main. error: %w", err)
	}

	_, err = tx.Exec(sizeQuery, id, position.Id, size.Dn, size.DnMm, size.Pn.Mpa, size.Pn.Kg, size.D4, size.D3, size.D2, size.D1,
		size.H, size.Another, size.UseDimensions,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query size. error: %w", err)
	}

	_, err = tx.Exec(materialQuery, id, position.Id, material.FillerId, material.FillerCode, material.TypeId, material.ConstructionId, material.ConstructionCode,
		material.RotaryPlugId, material.RotaryPlugCode, material.InnerRingId, material.InnerRindCode, material.OuterRingId, material.OuterRingCode,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query material. error: %w", err)
	}

	_, err = tx.Exec(designQuery, id, position.Id, design.HasJumper, design.JumperCode, design.JumperWidth, design.HasHole, design.HasRemovable,
		design.HasMounting, design.MountingCode, design.Drawing,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query design. error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to finish transaction. error: %w", err)
	}
	return nil
}

func (r *PositionPutgRepo) Delete(ctx context.Context, positionId string) error {
	mainQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionMainPutgTable)
	sizeQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionSizePutgTable)
	materialQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionMaterialPutgTable)
	designQuery := fmt.Sprintf(`DELETE FROM %s WHERE position_id=$1`, PositionDesignPutgTable)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction. error: %w", err)
	}

	_, err = tx.Exec(mainQuery, positionId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query main. error: %w", err)
	}
	_, err = tx.Exec(sizeQuery, positionId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query size. error: %w", err)
	}
	_, err = tx.Exec(materialQuery, positionId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query material. error: %w", err)
	}
	_, err = tx.Exec(designQuery, positionId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to complete query design. error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to finish transaction. error: %w", err)
	}
	return nil
}
