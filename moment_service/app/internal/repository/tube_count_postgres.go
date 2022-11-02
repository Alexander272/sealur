package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (r *DeviceRepo) GetTubeCount(ctx context.Context, req *device_api.GetTubeCountRequest) (tube []models.TubeCountDTO, err error) {
	query := fmt.Sprintf("SELECT id, value FROM %s", TubeCountTable)

	if err := r.db.Get(&tube, query); err != nil {
		return tube, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return tube, nil
}

func (r *DeviceRepo) CreateTubeCount(ctx context.Context, tube *device_api.CreateTubeCountRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (value) VALUES ($1) RETURNING id", TubeCountTable)

	row := r.db.QueryRow(query, tube.Value)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var returnedId int
	if err := row.Scan(&returnedId); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", returnedId), nil
}

func (r *DeviceRepo) CreateFewTubeCount(ctx context.Context, tube *device_api.CreateFewTubeCountRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (value) VALUES ", TubeCountTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(tube.TubeCount))

	for i, d := range tube.TubeCount {
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

func (r *DeviceRepo) UpdateTubeCount(ctx context.Context, tube *device_api.UpdateTubeCountRequest) error {
	query := fmt.Sprintf("UPDATE %s SET value=$1 WHERE id=$2", TubeCountTable)

	_, err := r.db.Exec(query, tube.Value, tube.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *DeviceRepo) DeleteTubeCount(ctx context.Context, tube *device_api.DeleteTubeCountRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TubeCountTable)

	if _, err := r.db.Exec(query, tube.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
