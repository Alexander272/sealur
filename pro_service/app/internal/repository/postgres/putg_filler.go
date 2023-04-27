package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_filler_api"
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
	query := fmt.Sprintf(`SELECT %s.id, %s.title as temperature, %s.title, description, designation
		FROM %s INNER JOIN %s ON temperature_id=%s.id WHERE construction_id=$1 ORDER BY code`,
		PutgFillerTable, TemperatureTable, PutgFillerTable, PutgFillerTable, TemperatureTable, TemperatureTable,
	)

	if err := r.db.Select(&data, query, req.ConstructionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pf := range data {
		fillers = append(fillers, &putg_filler_model.PutgFiller{
			Id:          pf.Id,
			Title:       pf.Title,
			Temperature: pf.Temperature,
			Description: pf.Description,
			Designation: pf.Designation,
		})
	}

	return fillers, nil
}

// TODO дописать оставшиеся функции
