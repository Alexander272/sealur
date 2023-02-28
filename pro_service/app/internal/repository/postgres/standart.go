package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/standard_model"
	"github.com/Alexander272/sealur_proto/api/pro/standard_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type StandardRepo struct {
	db *sqlx.DB
}

func NewStandardRepo(db *sqlx.DB) *StandardRepo {
	return &StandardRepo{db: db}
}

func (r *StandardRepo) GetAll(ctx context.Context, standard *standard_api.GetAllStandards) (standards []*standard_model.Standard, err error) {
	var data []models.Standard
	query := fmt.Sprintf("SELECT id, title, format FROM %s", StandardTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, s := range data {
		standards = append(standards, &standard_model.Standard{
			Id:     s.Id,
			Title:  s.Title,
			Format: s.Format,
		})
	}

	return standards, nil
}

func (r *StandardRepo) GetDefault(ctx context.Context) (standard *standard_model.Standard, err error) {
	var data models.Standard
	query := fmt.Sprintf("SELECT id, title, format FROM %s WHERE is_default=true LIMIT 1", StandardTable)

	if err := r.db.Get(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	standard = &standard_model.Standard{
		Id:     data.Id,
		Title:  data.Title,
		Format: data.Format,
	}

	return standard, nil
}

func (r *StandardRepo) Create(ctx context.Context, standard *standard_api.CreateStandard) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, format) VALUES ($1, $2, $3)", StandardTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, standard.Title, pq.Array(standard.Format))
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *StandardRepo) CreateSeveral(ctx context.Context, standards *standard_api.CreateSeveralStandard) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title, format) VALUES ", StandardTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(standards.Standards))

	c := 3
	for i, s := range standards.Standards {
		id := uuid.New()
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", i*c+1, i*c+2, i*c+3))
		args = append(args, id, s.Title, pq.Array(s.Format))
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *StandardRepo) Update(ctx context.Context, standard *standard_api.UpdateStandard) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, format=$2 WHERE id=$3", StandardTable)

	_, err := r.db.Exec(query, standard.Title, pq.Array(standard.Format), standard.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *StandardRepo) Delete(ctx context.Context, standard *standard_api.DeleteStandard) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", StandardTable)

	if _, err := r.db.Exec(query, standard.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
