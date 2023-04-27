package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_flange_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_flange_type_api"
	"github.com/jmoiron/sqlx"
)

type PutgFlangeTypeRepo struct {
	db *sqlx.DB
}

func NewPutgFlangeTypeRepo(db *sqlx.DB) *PutgFlangeTypeRepo {
	return &PutgFlangeTypeRepo{
		db: db,
	}
}

func (r *PutgFlangeTypeRepo) Get(ctx context.Context, req *putg_flange_type_api.GetPutgFlangeType,
) (flangeTypes []*putg_flange_type_model.PutgFlangeType, err error) {
	var data []models.PutgFlangeType
	query := fmt.Sprintf(`SELECT id, title, code FROM %s WHERE putg_standard_id=$1 ORDER BY code`, PutgFlangeTypeTable)

	if err := r.db.Select(&data, query, req.StandardId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pft := range data {
		flangeTypes = append(flangeTypes, &putg_flange_type_model.PutgFlangeType{
			Id:    pft.Id,
			Title: pft.Title,
			Code:  pft.Code,
		})
	}

	return flangeTypes, nil
}

// TODO дописать оставшиеся функции
