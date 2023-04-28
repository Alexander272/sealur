package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_type_api"
	"github.com/jmoiron/sqlx"
)

type PutgTypeRepo struct {
	db *sqlx.DB
}

func NewPutgTypeRepo(db *sqlx.DB) *PutgTypeRepo {
	return &PutgTypeRepo{
		db: db,
	}
}

func (r *PutgTypeRepo) Get(ctx context.Context, req *putg_type_api.GetPutgType) (types []*putg_type_model.PutgType, err error) {
	var data []models.PutgType
	query := fmt.Sprintf(`SELECT id, title, code FROM %s WHERE filler_id=$1 ORDER BY code`, PutgTypeTable)

	if err := r.db.Select(&data, query, req.BaseId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pt := range data {
		types = append(types, &putg_type_model.PutgType{
			Id:    pt.Id,
			Title: pt.Title,
			Code:  pt.Code,
		})
	}

	return types, nil
}
