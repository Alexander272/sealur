package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/rings_kit_type_model"
	"github.com/Alexander272/sealur_proto/api/pro/rings_kit_type_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RingsKitTypeRepo struct {
	db *sqlx.DB
}

func NewRingsKitTypeRepo(db *sqlx.DB) *RingsKitTypeRepo {
	return &RingsKitTypeRepo{
		db: db,
	}
}

type RingsKitType interface {
	GetAll(context.Context, *rings_kit_type_api.GetRingsKitTypes) ([]*rings_kit_type_model.RingsKitType, error)
	Create(context.Context, *rings_kit_type_api.CreateRingsKitType) error
	Update(context.Context, *rings_kit_type_api.UpdateRingsKitType) error
	Delete(context.Context, *rings_kit_type_api.DeleteRingsKitType) error
}

func (r *RingsKitTypeRepo) GetAll(ctx context.Context, req *rings_kit_type_api.GetRingsKitTypes) (kitTypes []*rings_kit_type_model.RingsKitType, err error) {
	var data []models.RingsKitType
	query := fmt.Sprintf(`SELECT id, code, title, description, image, designation
		FROM %s ORDER BY default_count`,
		RingsKitTypeTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, rt := range data {
		kitTypes = append(kitTypes, &rings_kit_type_model.RingsKitType{
			Id:          rt.Id,
			Code:        rt.Code,
			Title:       rt.Title,
			Description: rt.Description,
			Image:       rt.Image,
			Designation: rt.Designation,
		})
	}

	return kitTypes, nil
}

func (r *RingsKitTypeRepo) Create(ctx context.Context, kit *rings_kit_type_api.CreateRingsKitType) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, code, title, description, image, designation)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		RingsKitTypeTable,
	)
	id := uuid.New()

	_, err := r.db.Exec(query, id, kit.Code, kit.Title, kit.Description, kit.Image, kit.Designation)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingsKitTypeRepo) Update(ctx context.Context, kit *rings_kit_type_api.UpdateRingsKitType) error {
	query := fmt.Sprintf(`UPDATE %s SET code=$1, title=$2, description=$3, image=$4, designation=$5 WHERE id=$6`, RingsKitTypeTable)

	_, err := r.db.Exec(query, kit.Code, kit.Title, kit.Description, kit.Image, kit.Designation, kit.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingsKitTypeRepo) Delete(ctx context.Context, kit *rings_kit_type_api.DeleteRingsKitType) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RingsKitTypeTable)

	_, err := r.db.Exec(query, kit.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
