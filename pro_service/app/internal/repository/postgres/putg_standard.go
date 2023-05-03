package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_standard_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PutgStandardRepo struct {
	db *sqlx.DB
}

func NewPutgStandardRepo(db *sqlx.DB) *PutgStandardRepo {
	return &PutgStandardRepo{
		db: db,
	}
}

func (r *PutgStandardRepo) Get(ctx context.Context, req *putg_standard_api.GetPutgStandard) (standards []*putg_standard_model.PutgStandard, err error) {
	var data []models.PutgStandard
	query := fmt.Sprintf(`SELECT %s.id, title, code, dn_title, pn_title	FROM %s INNER JOIN %s ON flange_standard_id=%s.id ORDER BY count`,
		PutgStandardTable, PutgStandardTable, FlangeStandardTable, FlangeStandardTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, ps := range data {
		standards = append(standards, &putg_standard_model.PutgStandard{
			Id:      ps.Id,
			Title:   ps.Title,
			Code:    ps.Code,
			DnTitle: ps.DnTitle,
			PnTitle: ps.PnTitle,
		})
	}

	return standards, nil
}

// TODO можно еще создавать flange_standard если такого нету (на клиенте какой-нибудь select замутить с возможностью добавления новых позиций)
func (r *PutgStandardRepo) Create(ctx context.Context, st *putg_standard_api.CreatePutgStandard) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, flange_standard_id, count, dn_title, pn_title) VALUES ($1, $2, $3, $4, $5)`, PutgStandardTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, st.FlangeStandardId, st.Count, st.DnTitle, st.PnTitle)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgStandardRepo) Update(ctx context.Context, st *putg_standard_api.UpdatePutgStandard) error {
	query := fmt.Sprintf(`UPDATE %s SET flange_standard_id=$1, count=$2, dn_title=$3, pn_title=$4 WHERE id=$5`, PutgStandardTable)

	_, err := r.db.Exec(query, st.FlangeStandardId, st.Count, st.DnTitle, st.PnTitle, st.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgStandardRepo) Delete(ctx context.Context, st *putg_standard_api.DeletePutgStandard) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgStandardTable)

	_, err := r.db.Exec(query, st.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
