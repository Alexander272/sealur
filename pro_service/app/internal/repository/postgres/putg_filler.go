package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PutgFillerRepo struct {
	db *sqlx.DB
}

func NewPutgFillerRepo(db *sqlx.DB) *PutgFillerRepo {
	return &PutgFillerRepo{
		db: db,
	}
}

func (r *PutgFillerRepo) Get(ctx context.Context, req *putg_filler_api.GetPutgFiller) (fillers []*putg_filler_model.PutgFiller, err error) {
	var data []models.PutgFiller
	query := fmt.Sprintf(`SELECT %s.id, base_filler_id as base_id, %s.title as temperature, %s.title, description, designation
		FROM %s INNER JOIN %s ON base_filler_id=%s.id INNER JOIN %s ON temperature_id=%s.id WHERE putg_standard_id=$1 ORDER BY code`,
		PutgFillerTable, TemperatureTable, PutgFillerBaseTable, PutgFillerTable, PutgFillerBaseTable, PutgFillerBaseTable, TemperatureTable, TemperatureTable,
	)

	if err := r.db.Select(&data, query, req.StandardId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pf := range data {
		fillers = append(fillers, &putg_filler_model.PutgFiller{
			Id:          pf.Id,
			BaseId:      pf.BaseId,
			Title:       pf.Title,
			Temperature: pf.Temperature,
			Description: pf.Description,
			Designation: pf.Designation,
		})
	}

	return fillers, nil
}

func (r *PutgFillerRepo) Create(ctx context.Context, filler *putg_filler_api.CreatePutgFiller) error {
	id := uuid.New()
	query := fmt.Sprintf(`INSERT INTO %s(id, base_filler_id, putg_standard_id) VALUES ($1, $2, $3)`, PutgFillerTable)

	_, err := r.db.Exec(query, id, filler.FillerId, filler.StandardId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgFillerRepo) Update(ctx context.Context, filler *putg_filler_api.UpdatePutgFiller) error {
	query := fmt.Sprintf(`UPDATE %s SET base_filler_id=$1, putg_standard_id=$2 WHERE id=$3`, PutgFillerTable)

	_, err := r.db.Exec(query, filler.FillerId, filler.StandardId, filler.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgFillerRepo) Delete(ctx context.Context, filler *putg_filler_api.DeletePutgFiller) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgFillerTable)

	_, err := r.db.Exec(query, filler.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
