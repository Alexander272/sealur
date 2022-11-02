package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (r *DeviceRepo) GetTubeLenght(ctx context.Context, req *device_api.GetTubeLenghtRequest) (tube []models.TubeLenghtDTO, err error) {
	query := fmt.Sprintf("SELECT id, value FROM %s WHERE dev_id=$1", TubeLenghtTable)

	if err := r.db.Get(&tube, query, req.DevId); err != nil {
		return tube, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return tube, nil
}

func (r *DeviceRepo) CreateTubeLenght(ctx context.Context, tube *device_api.CreateTubeLenghtRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, value) VALUES ($1, $2) RETURNING id", TubeLenghtTable)

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

func (r *DeviceRepo) CreateFewTubeLenght(ctx context.Context, tube *device_api.CreateFewTubeLenghtRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, value) VALUES ", TubeLenghtTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(tube.TubeLenght))

	c := 2
	for i, d := range tube.TubeLenght {
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

func (r *DeviceRepo) UpdateTubeLenght(ctx context.Context, tube *device_api.UpdateTubeLenghtRequest) error {
	query := fmt.Sprintf("UPDATE %s SET dev_id=$1, value=$2 WHERE id=$3", TubeLenghtTable)

	_, err := r.db.Exec(query, tube.DevId, tube.Value, tube.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *DeviceRepo) DeleteTubeLenght(ctx context.Context, tube *device_api.DeleteTubeLenghtRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TubeLenghtTable)

	if _, err := r.db.Exec(query, tube.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
