package postgres

import (
	"context"
	"fmt"
	"math"

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
	query := fmt.Sprintf(`SELECT id, title, code, min_thickness, max_thickness FROM %s WHERE filler_id=$1 ORDER BY code`, PutgTypeTable)

	if err := r.db.Select(&data, query, req.BaseId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pt := range data {
		types = append(types, &putg_type_model.PutgType{
			Id:           pt.Id,
			Title:        pt.Title,
			Code:         pt.Code,
			MinThickness: math.Round(pt.MinThickness*1000) / 1000,
			MaxThickness: math.Round(pt.MaxThickness*1000) / 1000,
		})
	}

	return types, nil
}

func (r *PutgTypeRepo) Create(ctx context.Context, t *putg_type_api.CreatePutgType) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, title, code, filler_id, min_thickness, max_thickness) VALUES ($1, $2, $3, $4, $5, $6)`, PutgTypeTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, t.Title, t.Code, t.FillerId, t.MinThickness, t.MaxThickness)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgTypeRepo) Update(ctx context.Context, t *putg_type_api.UpdatePutgType) error {
	query := fmt.Sprintf(`UPDATE %s SET title=$1, code=$2, filler_id=$3, min_thickness=$4, max_thickness=$5 WHERE id=$6`, PutgTypeTable)

	_, err := r.db.Exec(query, t.Title, t.Code, t.FillerId, t.MinThickness, t.MaxThickness, t.Id)
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
