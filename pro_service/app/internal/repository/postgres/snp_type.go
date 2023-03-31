package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/snp_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/snp_api"
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

func (r *SNPTypeRepo) Get(ctx context.Context, req *snp_type_api.GetSnpTypes) (snp []*snp_type_model.SnpType, err error) {
	var data []models.SNPType
	query := fmt.Sprintf("SELECT id, title, code FROM %s WHERE snp_standard_id=$1", SnpTypeTable)

	if err := r.db.Select(&data, query, req.StandardId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, s := range data {
		snp = append(snp, &snp_type_model.SnpType{
			Id:    s.Id,
			Title: s.Title,
			Code:  s.Code,
		})
	}

	return snp, nil
}

func (r *SNPTypeRepo) GetWithFlange(ctx context.Context, req *snp_api.GetSnpData) (types []*snp_model.FlangeType, err error) {
	var data []models.SNPTypeWithFlange
	query := fmt.Sprintf(`SELECT %s.id as type_id, %s.title as type_title, %s.code as type_code, %s.id, %s.title, %s.code,
		%s.description, has_d4, has_d3, has_d2, has_d1
		FROM %s INNER JOIN %s ON %s.id=flange_type_id WHERE %s.snp_standard_id=$1 ORDER BY flange_type_id, %s.title`,
		SnpTypeTable, SnpTypeTable, SnpTypeTable, FlangeTypeTable, FlangeTypeTable, FlangeTypeTable, SnpTypeTable,
		SnpTypeTable, FlangeTypeTable, FlangeTypeTable, SnpTypeTable, SnpTypeTable,
	)

	if err := r.db.Select(&data, query, req.SnpStandardId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for i, s := range data {
		if i > 0 && s.Id == types[len(types)-1].Id {
			types[len(types)-1].Types = append(types[len(types)-1].Types, &snp_type_model.SnpType{
				Id:    s.TypeId,
				Code:  s.TypeCode,
				Title: s.TypeTitle,
				HasD4: s.HasD4,
				HasD3: s.HasD3,
				HasD2: s.HasD2,
				HasD1: s.HasD1,
				// Description: s.Description,
			})
		} else {
			types = append(types, &snp_model.FlangeType{
				Id:          s.Id,
				Title:       s.Title,
				Code:        s.Code,
				Description: s.Description,
				Types: []*snp_type_model.SnpType{{
					Id:    s.TypeId,
					Code:  s.TypeCode,
					Title: s.TypeTitle,
					HasD4: s.HasD4,
					HasD3: s.HasD3,
					HasD2: s.HasD2,
					HasD1: s.HasD1,
				}},
			})
		}
	}

	return types, nil
}

// TODO исправить запрос
func (r *SNPTypeRepo) Create(ctx context.Context, snp *snp_type_api.CreateSnpType) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, flange_type_id, code, snp_standard_id) VALUES ($1, $2, $3, $4, $5)", SnpTypeTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, snp.Title, snp.FlangeTypeId, snp.Code, snp.StandardId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPTypeRepo) CreateSeveral(ctx context.Context, snp *snp_type_api.CreateSeveralSnpType) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, code, flange_type_id, code, snp_standard_id) VALUES ", SnpTypeTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(snp.SnpTypes))

	c := 5
	for i, s := range snp.SnpTypes {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5))
		args = append(args, id, s.Title, s.FlangeTypeId, s.Code, s.StandardId)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPTypeRepo) Update(ctx context.Context, snp *snp_type_api.UpdateSnpType) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, flange_type_id=$2, code=$3, snp_standard_id=$4 WHERE id=$5", SnpTypeTable)

	_, err := r.db.Exec(query, snp.Title, snp.FlangeTypeId, snp.Code, snp.StandardId, snp.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPTypeRepo) Delete(ctx context.Context, snp *snp_type_api.DeleteSnpType) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", SnpTypeTable)

	if _, err := r.db.Exec(query, snp.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
