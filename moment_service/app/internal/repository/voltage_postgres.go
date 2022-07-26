package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (r *MaterialsRepo) CreateVoltage(ctx context.Context, voltage *moment_api.CreateVoltageRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (mark_id, temperature, voltage) VALUES ($1, $2, $3)", VoltageTable)

	args := make([]interface{}, 0)
	args = append(args, voltage.MarkId, voltage.Voltage[0].Temperature, voltage.Voltage[0].Voltage)

	for i, v := range voltage.Voltage {
		if i > 0 {
			query += fmt.Sprintf(", ($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
			args = append(args, voltage.MarkId, v.Temperature, v.Voltage)
		}
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialsRepo) UpdateVoltage(ctx context.Context, voltage *moment_api.UpdateVoltageRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	//TODO будут проблемы если нужно будет записать 0 в бд
	//? можно передавать инфинити если значения нет и делать проверку на равенство

	if voltage.MarkId != "" {
		setValues = append(setValues, fmt.Sprintf("mark_id=$%d", argId))
		args = append(args, voltage.MarkId)
		argId++
	}
	if voltage.Temperature != 0 {
		setValues = append(setValues, fmt.Sprintf("temperature=$%d", argId))
		args = append(args, voltage.Temperature)
		argId++
	}
	if voltage.Voltage != 0 {
		setValues = append(setValues, fmt.Sprintf("voltage=$%d", argId))
		args = append(args, voltage.Voltage)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", VoltageTable, setQuery, argId)

	args = append(args, voltage.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *MaterialsRepo) DeleteVoltage(ctx context.Context, voltage *moment_api.DeleteVoltageRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", VoltageTable)

	if _, err := r.db.Exec(query, voltage.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
