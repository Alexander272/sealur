package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/snp_type_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SNPTypeRepo struct {
	db *sqlx.DB
}

func NewSNPTypeRepo(db *sqlx.DB) *SNPTypeRepo {
	return &SNPTypeRepo{db: db}
}

func (r *SNPTypeRepo) Get(ctx context.Context, snp *snp_type_api.GetSnpTypes) (s []models.SNPType, err error) {
	query := fmt.Sprintf("SELECT id, title FROM %s WHERE flange_type_id=$1", SnpTypeTable)

	if err := r.db.Select(&s, query, snp.FlangeTypeId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return s, nil
}

func (r *SNPTypeRepo) Create(ctx context.Context, snp *snp_type_api.CreateSnpType) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, flange_type_id) VALUES ($1, $2, $3)", SnpTypeTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, snp.Title, snp.FlangeTypeId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPTypeRepo) CreateSeveral(ctx context.Context, snp *snp_type_api.CreateSeveralSnpType) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code, flange_type_id) VALUES ", SnpTypeTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(snp.SnpTypes))

	c := 3
	for i, s := range snp.SnpTypes {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", i*c+1, i*c+2, i*c+3))
		args = append(args, id, s.Title, s.FlangeTypeId)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPTypeRepo) Update(ctx context.Context, flange *snp_type_api.UpdateSnpType) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, flange_type_id=$2 WHERE id=$4", SnpTypeTable)

	_, err := r.db.Exec(query, flange.Title, flange.FlangeTypeId, flange.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPTypeRepo) Delete(ctx context.Context, flange *snp_type_api.DeleteSnpType) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", SnpTypeTable)

	if _, err := r.db.Exec(query, flange.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
