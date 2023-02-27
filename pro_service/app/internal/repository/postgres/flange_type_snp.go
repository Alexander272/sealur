package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/flange_type_snp_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/flange_type_snp_model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FlangeTypeRepo struct {
	db *sqlx.DB
}

func NewFlangeTypeRepo(db *sqlx.DB) *FlangeTypeRepo {
	return &FlangeTypeRepo{db: db}
}

func (r *FlangeTypeRepo) Get(ctx context.Context, flange *flange_type_snp_api.GetFlangeTypeSnp) (flanges []*flange_type_snp_model.FlangeTypeSnp, err error) {
	var data []models.FlangeTypeSnp
	query := fmt.Sprintf("SELECT id, title, code FROM %s WHERE standard_id=$1", FlangeTypeSNPTable)

	if err := r.db.Select(&data, query, flange.StandardId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, fts := range data {
		flanges = append(flanges, &flange_type_snp_model.FlangeTypeSnp{
			Id:    fts.Id,
			Title: fts.Title,
			Code:  fts.Code,
		})
	}

	return flanges, nil
}

func (r *FlangeTypeRepo) Create(ctx context.Context, flange *flange_type_snp_api.CreateFlangeTypeSnp) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code, standard_id) VALUES ($1, $2, $3, $4)", FlangeTypeSNPTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, flange.Title, flange.Code, flange.StandardId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeTypeRepo) CreateSeveral(ctx context.Context, flanges *flange_type_snp_api.CreateSeveralFlangeTypeSnp) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code, standard_id) VALUES ", FlangeTypeSNPTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(flanges.FlangeTypeSnp))

	c := 4
	for i, s := range flanges.FlangeTypeSnp {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4))
		args = append(args, id, s.Title, s.Code, s.StandardId)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeTypeRepo) Update(ctx context.Context, flange *flange_type_snp_api.UpdateFlangeTypeSnp) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, code=$2, standard_id=$3 WHERE id=$4", FlangeTypeSNPTable)

	_, err := r.db.Exec(query, flange.Title, flange.Code, flange.StandardId, flange.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeTypeRepo) Delete(ctx context.Context, flange *flange_type_snp_api.DeleteFlangeTypeSnp) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", FlangeTypeSNPTable)

	if _, err := r.db.Exec(query, flange.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
