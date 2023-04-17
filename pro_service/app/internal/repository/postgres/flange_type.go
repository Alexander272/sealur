package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_type_model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FlangeTypeRepo struct {
	db *sqlx.DB
}

func NewFlangeTypeRepo(db *sqlx.DB) *FlangeTypeRepo {
	return &FlangeTypeRepo{
		db: db,
	}
}

func (r *FlangeTypeRepo) Get(ctx context.Context, flange *flange_type_api.GetFlangeType) (flanges []*flange_type_model.FlangeType, err error) {
	var data []models.FlangeTypeSnp
	query := fmt.Sprintf("SELECT id, title, code FROM %s", FlangeTypeTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, fts := range data {
		flanges = append(flanges, &flange_type_model.FlangeType{
			Id:    fts.Id,
			Title: fts.Title,
			Code:  fts.Code,
		})
	}

	return flanges, nil
}

func (r *FlangeTypeRepo) Create(ctx context.Context, flange *flange_type_api.CreateFlangeType) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code) VALUES ($1, $2, $3)", FlangeTypeTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, flange.Title, flange.Code)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeTypeRepo) Update(ctx context.Context, flange *flange_type_api.UpdateFlangeType) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, code=$2 WHERE id=$3", FlangeTypeTable)

	_, err := r.db.Exec(query, flange.Title, flange.Code, flange.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeTypeRepo) Delete(ctx context.Context, flange *flange_type_api.DeleteFlangeType) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", FlangeTypeTable)

	if _, err := r.db.Exec(query, flange.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
