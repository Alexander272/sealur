package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (r *DeviceRepo) GetNumberOfMoves(ctx context.Context, req *device_api.GetNumberOfMovesRequest) (number []models.NumberOfMovesDTO, err error) {
	query := fmt.Sprintf("SELECT id, dev_id, count_id, value FROM %s", NumberOfMovesTable)
	var args []interface{}

	if req.DevId != "" {
		query += " WHERE dev_id=$1"
		args = append(args, req.DevId)
	}
	if req.CountId != "" {
		query += "AND count_id=$2"
		args = append(args, req.CountId)
	}

	if err := r.db.Select(&number, query, args...); err != nil {
		return number, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return number, nil
}

func (r *DeviceRepo) CreateNumberOfMoves(ctx context.Context, number *device_api.CreateNumberOfMovesRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, count_id, value) VALUES ($1, $2, $3) RETURNING id", NumberOfMovesTable)

	row := r.db.QueryRow(query, number.DevId, number.CountId, number.Value)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var returnedId int
	if err := row.Scan(&returnedId); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", returnedId), nil
}

func (r *DeviceRepo) CreateFewNumberOfMoves(ctx context.Context, number *device_api.CreateFewNumberOfMovesRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, count_id, value) VALUES ", NumberOfMovesTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(number.Number))

	c := 3
	for i, d := range number.Number {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", i*c+1, i*c+2, i*c+3))
		args = append(args, d.DevId, d.CountId, d.Value)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *DeviceRepo) UpdateNumberOfMoves(ctx context.Context, number *device_api.UpdateNumberOfMovesRequest) error {
	query := fmt.Sprintf("UPDATE %s SET dev_id=$1, count_id=$2, value=$3 WHERE id=$4", NumberOfMovesTable)

	_, err := r.db.Exec(query, number.DevId, number.CountId, number.Value, number.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *DeviceRepo) DeleteNumberOfMoves(ctx context.Context, number *device_api.DeleteNumberOfMovesRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", NumberOfMovesTable)

	if _, err := r.db.Exec(query, number.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
