package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_data_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_data_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SnpDataRepo struct {
	db *sqlx.DB
}

func NewSnpDataRepo(db *sqlx.DB) *SnpDataRepo {
	return &SnpDataRepo{db: db}
}

func (r *SnpDataRepo) Get(ctx context.Context, req *snp_data_api.GetSnpData) (snp *snp_data_model.SnpData, err error) {
	var data models.SnpData
	query := fmt.Sprintf(`SELECT id, has_inner_ring, has_frame, has_outer_ring, has_hole, has_jumper, has_mounting FROM %s 
		WHERE type_id=$1 LIMIT 1`, SnpDataTable)

	if err := r.db.Get(&data, query, req.TypeId); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	snp = &snp_data_model.SnpData{
		Id: data.Id,
		// TypeId:       sd.TypeId,
		HasInnerRing: data.HasInnerRing,
		HasFrame:     data.HasFrame,
		HasOuterRing: data.HasOuterRing,
		HasHole:      data.HasHole,
		HasJumper:    data.HasJumper,
		HasMounting:  data.HasMounting,
	}

	return snp, nil
}

func (r *SnpDataRepo) Create(ctx context.Context, snp *snp_data_api.CreateSnpData) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, type_id, standard_id, has_inner_ring, has_frame, has_outer_ring, has_hole, has_jumper, has_mounting)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, SnpDataTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, snp.TypeId, snp.StandardId, snp.HasInnerRing, snp.HasFrame, snp.HasOuterRing, snp.HasHole, snp.HasJumper, snp.HasMounting)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SnpDataRepo) Update(ctx context.Context, snp *snp_data_api.UpdateSnpData) error {
	query := fmt.Sprintf(`UPDATE %s	SET type_id=$1, standard_id=$2, has_inner_ring=$3, has_frame=$4, has_outer_ring=$5, 
		has_hole=$6, has_jumper=$7, has_mounting=$8 WHERE id=$9`, SnpDataTable)

	_, err := r.db.Exec(query, snp.TypeId, snp.StandardId, snp.HasInnerRing, snp.HasFrame, snp.HasOuterRing, snp.HasHole, snp.HasJumper, snp.HasMounting, snp.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SnpDataRepo) Delete(ctx context.Context, snp *snp_data_api.DeleteSnpData) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", SnpDataTable)

	if _, err := r.db.Exec(query, snp.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}