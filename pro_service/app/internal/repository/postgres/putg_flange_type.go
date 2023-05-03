package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_flange_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_flange_type_api"
	"github.com/google/uuid"
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

func (r *PutgFlangeTypeRepo) Create(ctx context.Context, fl *putg_flange_type_api.CreatePutgFlangeType) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, title, code, putg_standard_id) VALUES ($1, $2, $3, $4)`, PutgFlangeTypeTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, fl.Title, fl.Code, fl.StandardId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgFlangeTypeRepo) Update(ctx context.Context, fl *putg_flange_type_api.UpdatePutgFlangeType) error {
	query := fmt.Sprintf(`UPDATE %s SET title=$1, code=$2, putg_standard_id=$3 WHERE id=$4`, PutgFlangeTypeTable)

	_, err := r.db.Exec(query, fl.Title, fl.Code, fl.StandardId, fl.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgFlangeTypeRepo) Delete(ctx context.Context, fl *putg_flange_type_api.DeletePutgFlangeType) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgFlangeTypeTable)

	_, err := r.db.Exec(query, fl.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
