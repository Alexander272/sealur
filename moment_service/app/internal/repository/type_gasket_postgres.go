package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (r *GasketRepo) GetTypeGasket(ctx context.Context, req *moment_api.GetGasketTypeRequest) (types []models.TypeGasketDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title, label FROM %s ORDER BY id`, TypeGasketTable)

	if err := r.db.Select(&types, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return types, nil
}

func (r *GasketRepo) CreateTypeGasket(ctx context.Context, typeGasket *moment_api.CreateGasketTypeRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title, label) VALUES ($1, $2) RETURNING id", TypeGasketTable)

	row := r.db.QueryRow(query, typeGasket.Title, typeGasket.Label)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *GasketRepo) UpdateTypeGasket(ctx context.Context, typeGasket *moment_api.UpdateGasketTypeRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if typeGasket.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, typeGasket.Title)
		argId++
	}
	if typeGasket.Label != "" {
		setValues = append(setValues, fmt.Sprintf("label=$%d", argId))
		args = append(args, typeGasket.Label)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", TypeGasketTable, setQuery, argId)

	args = append(args, typeGasket.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *GasketRepo) DeleteTypeGasket(ctx context.Context, typeGasket *moment_api.DeleteGasketTypeRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TypeGasketTable)

	if _, err := r.db.Exec(query, typeGasket.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
