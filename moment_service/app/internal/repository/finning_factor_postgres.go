package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (r *DeviceRepo) GetFinningFactor(ctx context.Context, req *device_api.GetFinningFactorRequest) (factor []models.FinningFactorDTO, err error) {
	query := fmt.Sprintf("SELECT id, value FROM %s WHERE dev_id=$1", FinningFactorTable)

	if err := r.db.Select(&factor, query, req.DevId); err != nil {
		return factor, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return factor, nil
}

func (r *DeviceRepo) CreateFinningFactor(ctx context.Context, factor *device_api.CreateFinningFactorRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, value) VALUES ($1, $2) RETURNING id", FinningFactorTable)

	row := r.db.QueryRow(query, factor.DevId, factor.Value)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var returnedId int
	if err := row.Scan(&returnedId); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", returnedId), nil
}

func (r *DeviceRepo) CreateFewFinningFactor(ctx context.Context, factor *device_api.CreateFewFinningFactorRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, value) VALUES ", FinningFactorTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(factor.FinnigFactor))

	c := 2
	for i, d := range factor.FinnigFactor {
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

func (r *DeviceRepo) UpdateFinningFactor(ctx context.Context, factor *device_api.UpdateFinningFactorRequest) error {
	query := fmt.Sprintf("UPDATE %s SET dev_id=$1, value=$2 WHERE id=$3", FinningFactorTable)

	_, err := r.db.Exec(query, factor.DevId, factor.Value, factor.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *DeviceRepo) DeleteFinningFactor(ctx context.Context, factor *device_api.DeleteFinningFactorRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", FinningFactorTable)

	if _, err := r.db.Exec(query, factor.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
