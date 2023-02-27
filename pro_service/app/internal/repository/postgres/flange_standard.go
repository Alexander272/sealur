package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/flange_standard_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_standard_model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FlangeStandardRepo struct {
	db *sqlx.DB
}

func NewFlangeStandardRepo(db *sqlx.DB) *FlangeStandardRepo {
	return &FlangeStandardRepo{db: db}
}

func (r *FlangeStandardRepo) GetAll(ctx context.Context, flange *flange_standard_api.GetAllFlangeStandards) (flanges []*flange_standard_model.FlangeStandard, err error) {
	var data []models.FlangeStandard
	query := fmt.Sprintf("SELECT id, title, code FROM %s", FlangeStandardTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, fs := range data {
		flanges = append(flanges, &flange_standard_model.FlangeStandard{
			Id:    fs.Id,
			Title: fs.Title,
			Code:  fs.Code,
		})
	}

	return flanges, nil
}

func (r *FlangeStandardRepo) Create(ctx context.Context, flange *flange_standard_api.CreateFlangeStandard) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code) VALUES ($1, $2, $3)", FlangeStandardTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, flange.Title, flange.Code)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeStandardRepo) CreateSeveral(ctx context.Context, flanges *flange_standard_api.CreateSeveralFlangeStandard) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code) VALUES ", FlangeStandardTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(flanges.FlangeStandards))

	c := 3
	for i, s := range flanges.FlangeStandards {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", i*c+1, i*c+2, i*c+3))
		args = append(args, id, s.Title, s.Code)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeStandardRepo) Update(ctx context.Context, flange *flange_standard_api.UpdateFlangeStandard) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, code=$2 WHERE id=$3", FlangeStandardTable)

	_, err := r.db.Exec(query, flange.Title, flange.Code, flange.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeStandardRepo) Delete(ctx context.Context, flange *flange_standard_api.DeleteFlangeStandard) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", FlangeStandardTable)

	if _, err := r.db.Exec(query, flange.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
