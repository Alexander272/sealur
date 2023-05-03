package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_data_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_data_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PutgDataRepo struct {
	db *sqlx.DB
}

func NewPutgDataRepo(db *sqlx.DB) *PutgDataRepo {
	return &PutgDataRepo{
		db: db,
	}
}

func (r *PutgDataRepo) Get(ctx context.Context, req *putg_data_api.GetPutgData) (putg *putg_data_model.PutgData, err error) {
	var data models.PutgData
	query := fmt.Sprintf(`SELECT id, filler_id, has_jumper, has_hole, has_removable, has_mounting, has_coating FROM %s WHERE filler_id=$1`, PutgDataTable)

	if err := r.db.Get(&data, query, req.FillerId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	putg = &putg_data_model.PutgData{
		Id:           data.Id,
		HasJumper:    data.HasJumper,
		HasHole:      data.HasHole,
		HasRemovable: data.HasRemovable,
		HasMounting:  data.HasMounting,
		HasCoating:   data.HasCoating,
	}

	return putg, nil
}

func (r *PutgDataRepo) GetByConstruction(ctx context.Context, req *putg_data_api.GetPutgData) (putg []*putg_data_model.PutgData, err error) {
	var data []models.PutgData
	query := fmt.Sprintf(`SELECT %s.id, filler_id, has_jumper, has_hole, has_removable, has_mounting, has_coating FROM %s 
		INNER JOIN %s ON filler_id=%s.id WHERE construction_id=$1 ORDER BY code`,
		PutgDataTable, PutgDataTable, PutgFillerTable, PutgFillerTable,
	)

	if err := r.db.Select(&data, query, req.ConstructionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pd := range data {
		putg = append(putg, &putg_data_model.PutgData{
			Id:           pd.Id,
			FillerId:     pd.FillerId,
			HasJumper:    pd.HasJumper,
			HasHole:      pd.HasHole,
			HasRemovable: pd.HasRemovable,
			HasMounting:  pd.HasMounting,
			HasCoating:   pd.HasCoating,
		})
	}

	return putg, nil
}

func (r *PutgDataRepo) Create(ctx context.Context, data *putg_data_api.CreatePutgData) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, filler_id, has_jumper, has_hole, has_removable, has_mounting, has_coating)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, PutgDataTable,
	)
	id := uuid.New()

	_, err := r.db.Exec(query, id, data.FillerId, data.HasJumper, data.HasHole, data.HasRemovable, data.HasMounting, data.HasCoating)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgDataRepo) Update(ctx context.Context, data *putg_data_api.UpdatePutgData) error {
	query := fmt.Sprintf(`UPDATE %s SET filler_id=$1, has_jumper=$2, has_hole=$3, has_removable=$4, has_mounting=$5, has_coating=$6 
		WHERE id=$7`, PutgDataTable,
	)

	_, err := r.db.Exec(query, data.FillerId, data.HasJumper, data.HasHole, data.HasRemovable, data.HasMounting, data.HasCoating, data.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgDataRepo) Delete(ctx context.Context, data *putg_data_api.DeletePutgData) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgDataTable)

	_, err := r.db.Exec(query, data.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
