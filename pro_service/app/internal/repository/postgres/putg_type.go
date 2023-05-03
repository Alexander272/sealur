package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_type_api"
	"github.com/google/uuid"
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

func (r *PutgTypeRepo) Create(ctx context.Context, t *putg_type_api.CreatePutgType) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, title, code, filler_id) VALUES ($1, $2, $3, $4)`, PutgTypeTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, t.Title, t.Code, t.FillerId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgTypeRepo) Update(ctx context.Context, t *putg_type_api.UpdatePutgType) error {
	query := fmt.Sprintf(`UPDATE %s SET title=$1, code=$2, filler_id=$3 WHERE id=$4`, PutgTypeTable)

	_, err := r.db.Exec(query, t.Title, t.Code, t.FillerId, t.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgTypeRepo) Delete(ctx context.Context, t *putg_type_api.DeletePutgType) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgTypeTable)

	_, err := r.db.Exec(query, t.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
