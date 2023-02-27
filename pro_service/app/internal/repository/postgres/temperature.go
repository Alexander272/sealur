package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/temperature_model"
	"github.com/Alexander272/sealur_proto/api/pro/temperature_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TemperatureRepo struct {
	db *sqlx.DB
}

func NewTemperatureRepo(db *sqlx.DB) *TemperatureRepo {
	return &TemperatureRepo{db: db}
}

func (r *TemperatureRepo) GetAll(ctx context.Context, temp *temperature_api.GetAllTemperatures) (temperatures []*temperature_model.Temperature, err error) {
	var data []models.Temperature
	query := fmt.Sprintf("SELECT id, title FROM %s", TemperatureTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, t := range data {
		temperatures = append(temperatures, &temperature_model.Temperature{
			Id:    t.Id,
			Title: t.Title,
		})
	}

	return temperatures, nil
}

func (r *TemperatureRepo) Create(ctx context.Context, temperature *temperature_api.CreateTemperature) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title) VALUES ($1, $2)", TemperatureTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, temperature.Title)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *TemperatureRepo) CreateSeveral(ctx context.Context, temperatures *temperature_api.CreateSeveralTemperature) error {
	query := fmt.Sprintf("INSERT INTO %s (id, title) VALUES ", TemperatureTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(temperatures.Temperatures))

	c := 2
	for i, m := range temperatures.Temperatures {
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

func (r *TemperatureRepo) Update(ctx context.Context, temperature *temperature_api.UpdateTemperature) error {
	query := fmt.Sprintf("UPDATE %s	SET title=$1, code=$2 WHERE id=$3", TemperatureTable)

	_, err := r.db.Exec(query, temperature.Title, temperature.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *TemperatureRepo) Delete(ctx context.Context, temperature *temperature_api.DeleteTemperature) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TemperatureTable)

	if _, err := r.db.Exec(query, temperature.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
