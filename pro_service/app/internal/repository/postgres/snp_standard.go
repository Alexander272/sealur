package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_standard_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SnpStandardRepo struct {
	db *sqlx.DB
}

func NewSnpStandardRepo(db *sqlx.DB) *SnpStandardRepo {
	return &SnpStandardRepo{db: db}
}

func (r *SnpStandardRepo) GetAll(ctx context.Context, req *snp_standard_api.GetAllSnpStandards) (standards []*snp_standard_model.SnpStandard, err error) {
	var data []models.SnpStandard
	query := fmt.Sprintf(`SELECT %s.id, dn_title, pn_title, %s.title as standard_title, %s.format as standard_format, %s.title as flange_title, %s.code as flange_code
		FROM %s INNER JOIN %s ON %s.id=standard_id INNER JOIN %s ON %s.id=flange_standard_id`,
		SnpStandardTable, StandardTable, StandardTable, FlangeStandardTable, FlangeStandardTable,
		SnpStandardTable, StandardTable, StandardTable, FlangeStandardTable, FlangeStandardTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	//TODO добавить все поля
	for _, s := range data {
		standards = append(standards, &snp_standard_model.SnpStandard{
			Id:      s.Id,
			DnTitle: s.DnTitle,
			PnTitle: s.PnTitle,
		})
	}

	return standards, nil
}

func (r *SnpStandardRepo) Create(ctx context.Context, standard *snp_standard_api.CreateSnpStandard) error {
	query := fmt.Sprintf("INSERT INTO %s(id, standard_id, flange_standard_id, dn_title, pn_title) VALUES ($1, $2, $3, $4, $5)", SnpStandardTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, standard.StandardId, standard.FlangeStandardId, standard.DnTitle, standard.PnTitle)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SnpStandardRepo) CreateSeveral(ctx context.Context, standard *snp_standard_api.CreateSeveralSnpStandard) error {
	query := fmt.Sprintf("INSERT INTO %s (id, standard_id, flange_standard_id, dn_title, pn_title) VALUES ", SnpStandardTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(standard.SnpStandards))

	c := 5
	for i, f := range standard.SnpStandards {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5))
		args = append(args, id, f.StandardId, f.FlangeStandardId, f.DnTitle, f.PnTitle)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SnpStandardRepo) Update(ctx context.Context, standard *snp_standard_api.UpdateSnpStandard) error {
	query := fmt.Sprintf("UPDATE %s	SET standard_id=$1, flange_standard_id=$2, dn_title=$3, pn_title=$4 WHERE id=$5", SnpStandardTable)

	_, err := r.db.Exec(query, standard.StandardId, standard.FlangeStandardId, standard.DnTitle, standard.PnTitle, standard.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SnpStandardRepo) Delete(ctx context.Context, standard *snp_standard_api.DeleteSnpStandard) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", SnpStandardTable)

	if _, err := r.db.Exec(query, standard.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
