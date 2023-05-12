package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_base_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PutgBaseFillerRepo struct {
	db *sqlx.DB
}

func NewPutgBaseFillerRepo(db *sqlx.DB) *PutgBaseFillerRepo {
	return &PutgBaseFillerRepo{
		db: db,
	}
}

func (r *PutgBaseFillerRepo) Get(ctx context.Context, req *putg_filler_base_api.GetPutgBaseFiller) (fillers []*putg_filler_model.PutgFiller, err error) {
	var data []models.PutgFiller
	query := fmt.Sprintf(`SELECT %s.id, %s.title as temperature, %s.title, description, designation
	 	FROM %s INNER JOIN %s ON temperature_id=%s.id ORDER BY code`,
		PutgFillerBaseTable, TemperatureTable, PutgFillerBaseTable, PutgFillerBaseTable, TemperatureTable, TemperatureTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pf := range data {
		fillers = append(fillers, &putg_filler_model.PutgFiller{
			Id:          pf.Id,
			Temperature: pf.Temperature,
			Title:       pf.Title,
			Description: pf.Description,
			Designation: pf.Designation,
		})
	}

	return fillers, nil
}

func (r *PutgBaseFillerRepo) Create(ctx context.Context, filler *putg_filler_base_api.CreatePutgBaseFiller) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, temperature_id, title, description, designation, code)
		VALUES ($1, $2, $3, $4, $5, $6)`, PutgFillerBaseTable,
	)
	id := uuid.New()

	_, err := r.db.Exec(query, id, filler.TemperatureId, filler.Title, filler.Description, filler.Designation, filler.Code)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgBaseFillerRepo) Update(ctx context.Context, filler *putg_filler_base_api.UpdatePutgBaseFiller) error {
	query := fmt.Sprintf(`UPDATE %s SET temperature_id=$1, title=$2, description=$3, designation=$4, code=$5 WHERE id=$6`, PutgFillerBaseTable)

	_, err := r.db.Exec(query, filler.TemperatureId, filler.Title, filler.Description, filler.Designation, filler.Code, filler.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgBaseFillerRepo) Delete(ctx context.Context, filler *putg_filler_base_api.DeletePutgBaseFiller) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgFillerBaseTable)

	_, err := r.db.Exec(query, filler.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
