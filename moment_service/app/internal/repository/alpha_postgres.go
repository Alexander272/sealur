package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (r *MaterialsRepo) CreateAlpha(ctx context.Context, alpha *moment_api.CreateAlphaRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (mark_id, temperature, alpha) VALUES ($1, $2, $3)", AlphaTable)

	args := make([]interface{}, 0)
	args = append(args, alpha.MarkId, alpha.Alpha[0].Temperature, alpha.Alpha[0].Alpha)

	for i, a := range alpha.Alpha {
		if i > 0 {
			query += fmt.Sprintf(", ($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
			args = append(args, alpha.MarkId, a.Temperature, a.Alpha)
		}
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialsRepo) UpateAlpha(ctx context.Context, alpha *moment_api.UpdateAlphaRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if alpha.MarkId != "" {
		setValues = append(setValues, fmt.Sprintf("mark_id=$%d", argId))
		args = append(args, alpha.MarkId)
		argId++
	}
	if alpha.Temperature != 0 {
		setValues = append(setValues, fmt.Sprintf("temperature=$%d", argId))
		args = append(args, alpha.Temperature)
		argId++
	}
	if alpha.Alpha != 0 {
		setValues = append(setValues, fmt.Sprintf("alpha=$%d", argId))
		args = append(args, alpha.Alpha)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", AlphaTable, setQuery, argId)

	args = append(args, alpha.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *MaterialsRepo) DeleteAlpha(ctx context.Context, alpha *moment_api.DeleteAlphaRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", AlphaTable)

	if _, err := r.db.Exec(query, alpha.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
