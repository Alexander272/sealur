package repository

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/material_api"
)

func (r *MaterialsRepo) CreateVoltage(ctx context.Context, voltage *material_api.CreateVoltageRequest) error {
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

func (r *MaterialsRepo) UpdateVoltage(ctx context.Context, voltage *material_api.UpdateVoltageRequest) error {
	query := fmt.Sprintf("UPDATE %s SET temperature=$1, voltage=$2 WHERE id=$3", VoltageTable)

	_, err := r.db.Exec(query, voltage.Temperature, voltage.Voltage, voltage.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *MaterialsRepo) DeleteVoltage(ctx context.Context, voltage *material_api.DeleteVoltageRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", VoltageTable)

	if _, err := r.db.Exec(query, voltage.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
