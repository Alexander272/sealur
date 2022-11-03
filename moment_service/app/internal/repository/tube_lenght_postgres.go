package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (r *DeviceRepo) GetTubeLength(ctx context.Context, req *device_api.GetTubeLengthRequest) (tube []models.TubeLengthDTO, err error) {
	query := fmt.Sprintf("SELECT id, value FROM %s WHERE dev_id=$1", TubeLengthTable)

	if err := r.db.Get(&tube, query, req.DevId); err != nil {
		return tube, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return tube, nil
}

func (r *DeviceRepo) CreateTubeLength(ctx context.Context, tube *device_api.CreateTubeLengthRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, value) VALUES ($1, $2) RETURNING id", TubeLengthTable)

	row := r.db.QueryRow(query, tube.DevId, tube.Value)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var returnedId int
	if err := row.Scan(&returnedId); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", returnedId), nil
}

func (r *DeviceRepo) CreateFewTubeLength(ctx context.Context, tube *device_api.CreateFewTubeLengthRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, value) VALUES ", TubeLengthTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(tube.TubeLength))

	c := 2
	for i, d := range tube.TubeLength {
		values = append(values, fmt.Sprintf("($%d, $%d)", i*c+1, i*c+2))
		args = append(args, d.DevId, d.Value)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *DeviceRepo) UpdateTubeLength(ctx context.Context, tube *device_api.UpdateTubeLengthRequest) error {
	query := fmt.Sprintf("UPDATE %s SET dev_id=$1, value=$2 WHERE id=$3", TubeLengthTable)

	_, err := r.db.Exec(query, tube.DevId, tube.Value, tube.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *DeviceRepo) DeleteTubeLength(ctx context.Context, tube *device_api.DeleteTubeLengthRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TubeLengthTable)

	if _, err := r.db.Exec(query, tube.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
