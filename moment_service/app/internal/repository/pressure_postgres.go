package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (r *DeviceRepo) GetPressure(ctx context.Context, req *device_api.GetPressureRequest) (pressure []models.PressureDTO, err error) {
	query := fmt.Sprintf("SELECT id, value FROM %s", PressureTable)

	if err := r.db.Get(&pressure, query); err != nil {
		return pressure, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return pressure, nil
}

func (r *DeviceRepo) CreatePressure(ctx context.Context, pres *device_api.CreatePressureRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (value) VALUES ($1) RETURNING id", PressureTable)

	row := r.db.QueryRow(query, pres.Value)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var returnedId int
	if err := row.Scan(&returnedId); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", returnedId), nil
}

func (r *DeviceRepo) CreateFewPressure(ctx context.Context, pres *device_api.CreateFewPressureRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (value) VALUES ", PressureTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(pres.Pressure))

	for i, d := range pres.Pressure {
		values = append(values, fmt.Sprintf("($%d)", i+1))
		args = append(args, d.Value)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *DeviceRepo) UpdatePressure(ctx context.Context, pres *device_api.UpdatePressureRequest) error {
	query := fmt.Sprintf("UPDATE %s SET value=$1 WHERE id=$2", PressureTable)

	_, err := r.db.Exec(query, pres.Value, pres.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *DeviceRepo) DeletePressure(ctx context.Context, pres *device_api.DeletePressureRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", PressureTable)

	if _, err := r.db.Exec(query, pres.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
