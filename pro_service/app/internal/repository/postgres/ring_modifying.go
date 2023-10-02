package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/ring_modifying_model"
	"github.com/Alexander272/sealur_proto/api/pro/ring_modifying_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RingModifyingRepo struct {
	db *sqlx.DB
}

func NewRingModifyingRepo(db *sqlx.DB) *RingModifyingRepo {
	return &RingModifyingRepo{
		db: db,
	}
}

type RingModifying interface {
	GetAll(context.Context, *ring_modifying_api.GetRingModifying) ([]*ring_modifying_model.RingModifying, error)
	Create(context.Context, *ring_modifying_api.CreateRingModifying) error
	Update(context.Context, *ring_modifying_api.UpdateRingModifying) error
	Delete(context.Context, *ring_modifying_api.DeleteRingModifying) error
}

func (r *RingModifyingRepo) GetAll(ctx context.Context, req *ring_modifying_api.GetRingModifying) (m []*ring_modifying_model.RingModifying, err error) {
	var data []models.RingModifying
	query := fmt.Sprintf(`SELECT id, code, title, description, designation FROM %s WHERE is_show=true`, RingModifyingTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, rt := range data {
		m = append(m, &ring_modifying_model.RingModifying{
			Id:          rt.Id,
			Code:        rt.Code,
			Title:       rt.Title,
			Description: rt.Description,
			Designation: rt.Designation,
		})
	}

	return m, nil
}

func (r *RingModifyingRepo) Create(ctx context.Context, m *ring_modifying_api.CreateRingModifying) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, code, description, designation) VALUES ($1, $2, $3, $4)`, RingModifyingTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, m.Code, m.Description, m.Designation)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingModifyingRepo) Update(ctx context.Context, m *ring_modifying_api.UpdateRingModifying) error {
	query := fmt.Sprintf(`UPDATE %s SET code=$1, description=$2, designation=$3 WHERE id=$4`, RingModifyingTable)

	_, err := r.db.Exec(query, m.Code, m.Description, m.Designation, m.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RingModifyingRepo) Delete(ctx context.Context, m *ring_modifying_api.DeleteRingModifying) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RingModifyingTable)

	_, err := r.db.Exec(query, m.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
