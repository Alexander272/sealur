package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/snp_filler_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SNPFillerRepo struct {
	db *sqlx.DB
}

func NewSNPFillerRepo(db *sqlx.DB) *SNPFillerRepo {
	return &SNPFillerRepo{db: db}
}

func (r *SNPFillerRepo) GetAll(ctx context.Context, fil *snp_filler_api.GetSnpFillers) (fillers []models.SNPFiller, err error) {
	query := fmt.Sprintf("SELECT id, title, code, description FROM %s", SnpFillerTable)

	if err := r.db.Select(&fillers, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return fillers, nil
}

func (r *SNPFillerRepo) Create(ctx context.Context, filler *snp_filler_api.CreateSnpFiller) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code, description) VALUES ($1, $2, $3, $4)", SnpFillerTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, filler.Title, filler.Code, filler.Description)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPFillerRepo) CreateSeveral(ctx context.Context, fillers *snp_filler_api.CreateSeveralSnpFiller) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code, description) VALUES ", SnpFillerTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(fillers.SnpFillers))

	c := 4
	for i, f := range fillers.SnpFillers {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4))
		args = append(args, id, f.Title, f.Code, f.Description)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPFillerRepo) Update(ctx context.Context, filler *snp_filler_api.UpdateSnpFiller) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, code=$2, description=$3 WHERE id=$4", SnpFillerTable)

	_, err := r.db.Exec(query, filler.Title, filler.Code, filler.Description, filler.Id)
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
