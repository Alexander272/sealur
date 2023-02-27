package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/mounting_model"
	"github.com/Alexander272/sealur_proto/api/pro/mounting_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MountingRepo struct {
	db *sqlx.DB
}

func NewMountingRepo(db *sqlx.DB) *MountingRepo {
	return &MountingRepo{db: db}
}

func (r *MountingRepo) GetAll(ctx context.Context, mount *mounting_api.GetAllMountings) (mounting []*mounting_model.Mounting, err error) {
	var data []models.Mounting
	query := fmt.Sprintf("SELECT id, title FROM %s", MountingTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, m := range data {
		mounting = append(mounting, &mounting_model.Mounting{
			Id:    m.Id,
			Title: m.Title,
		})
	}

	return mounting, nil
}

func (r *MountingRepo) Create(ctx context.Context, mounting *mounting_api.CreateMounting) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title) VALUES ($1, $2)", MountingTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, mounting.Title)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MountingRepo) CreateSeveral(ctx context.Context, mounting *mounting_api.CreateSeveralMounting) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title) VALUES ", MountingTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(mounting.Mounting))

	c := 2
	for i, m := range mounting.Mounting {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d)", i*c+1, i*c+2))
		args = append(args, id, m.Title)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MountingRepo) Update(ctx context.Context, mounting *mounting_api.UpdateMounting) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, code=$2 WHERE id=$3", MountingTable)

	_, err := r.db.Exec(query, mounting.Title, mounting.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MountingRepo) Delete(ctx context.Context, mounting *mounting_api.DeleteMounting) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", MountingTable)

	if _, err := r.db.Exec(query, mounting.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
