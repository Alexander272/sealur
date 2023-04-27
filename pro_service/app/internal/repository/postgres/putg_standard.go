package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_standard_api"
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

// TODO дописать оставшиеся функции
