package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_filler_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SNPFillerRepo struct {
	db *sqlx.DB
}

func NewSNPFillerRepo(db *sqlx.DB) *SNPFillerRepo {
	return &SNPFillerRepo{db: db}
}

func (r *SNPFillerRepo) GetAll(ctx context.Context, req *snp_filler_api.GetSnpFillers) (fillers []*snp_filler_model.SnpFiller, err error) {
	var data []models.SnpFillerNew
	query := fmt.Sprintf(`SELECT %s.id, %s.title, base_code, code, description, designation, disabled_types, %s.title as temperature
		FROM %s INNER JOIN %s on %s.id=temperature_id WHERE standard_id=$1 ORDER BY base_code`,
		SnpFillerNewTable, SnpFillerNewTable, TemperatureTable, SnpFillerNewTable, TemperatureTable, TemperatureTable,
	)

	if err := r.db.Select(&data, query, req.StandardId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, s := range data {
		fillers = append(fillers, &snp_filler_model.SnpFiller{
			Id:            s.Id,
			Temperature:   s.Temperature,
			BaseCode:      s.BaseCode,
			Code:          s.Code,
			Title:         s.Title,
			Description:   s.Description,
			Designation:   s.Designation,
			DisabledTypes: s.DisabledTypes,
		})
	}

	return fillers, nil
}

func (r *SNPFillerRepo) Create(ctx context.Context, filler *snp_filler_api.CreateSnpFiller) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, standard_id, temperature_id, base_code, code, title, description, designation, disabled_types) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, SnpFillerTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, filler.StandardId, filler.TemperatureId, filler.BaseCode, filler.Code, filler.Title, filler.Description,
		filler.Designation, pq.Array(filler.DisabledTypes),
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPFillerRepo) CreateSeveral(ctx context.Context, fillers *snp_filler_api.CreateSeveralSnpFiller) error {
	query := fmt.Sprintf("INSERT INTO %s (id, standard_id, temperature_id, base_code, code, title, description, designation, disabled_types) VALUES ", SnpFillerTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(fillers.SnpFillers))

	c := 9
	for i, f := range fillers.SnpFillers {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5, i*c+6, i*c+7, i*c+8, i*c+9))
		args = append(args, id, f.StandardId, f.TemperatureId, f.BaseCode, f.Code, f.Title, f.Description, f.Designation, pq.Array(f.DisabledTypes))
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPFillerRepo) Update(ctx context.Context, filler *snp_filler_api.UpdateSnpFiller) error {
	query := fmt.Sprintf(`UPDATE %s	SET standard_id=$1, temperature_id=$2, base_code=$3, code=$4, title=$5, description=$6, designation=$7,
		disabled_types=$8 WHERE id=$9`, SnpFillerTable)

	_, err := r.db.Exec(query, filler.StandardId, filler.TemperatureId, filler.BaseCode, filler.Code, filler.Title, filler.Description,
		filler.Designation, pq.Array(filler.DisabledTypes), filler.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPFillerRepo) Delete(ctx context.Context, filler *snp_filler_api.DeleteSnpFiller) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", SnpFillerTable)

	if _, err := r.db.Exec(query, filler.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
