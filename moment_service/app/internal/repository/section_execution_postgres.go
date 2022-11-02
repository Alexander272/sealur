package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (r *DeviceRepo) GetSectionExecution(ctx context.Context, req *device_api.GetSectionExecutionRequest) (section []models.SectionExecutionDTO, err error) {
	query := fmt.Sprintf("SELECT id, value FROM %s WHERE dev_id=$1", SectionExecutionTable)

	if err := r.db.Get(&section, query, req.DevId); err != nil {
		return section, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return section, nil
}

func (r *DeviceRepo) CreateSectionExecution(ctx context.Context, section *device_api.CreateSectionExecutionRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, value) VALUES ($1, $2) RETURNING id", SectionExecutionTable)

	row := r.db.QueryRow(query, section.DevId, section.Value)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var returnedId int
	if err := row.Scan(&returnedId); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", returnedId), nil
}

func (r *DeviceRepo) CreateFewSectionExecution(ctx context.Context, section *device_api.CreateFewSectionExecutionRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (dev_id, value) VALUES ", SectionExecutionTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(section.Section))

	c := 2
	for i, d := range section.Section {
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

func (r *DeviceRepo) UpdateSectionExecution(ctx context.Context, section *device_api.UpdateSectionExecutionRequest) error {
	query := fmt.Sprintf("UPDATE %s SET dev_id=$1, value=$2 WHERE id=$3", SectionExecutionTable)

	_, err := r.db.Exec(query, section.DevId, section.Value, section.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *DeviceRepo) DeleteSectionExecution(ctx context.Context, section *device_api.DeleteSectionExecutionRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", SectionExecutionTable)

	if _, err := r.db.Exec(query, section.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
