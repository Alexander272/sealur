package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/gasket_api"
)

func (r *GasketRepo) GetEnv(ctx context.Context, req *gasket_api.GetEnvRequest) (env []models.EnvDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title FROM %s ORDER BY id`, EnvTable)

	if err := r.db.Select(&env, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return env, nil
}

func (r *GasketRepo) CreateEnv(ctx context.Context, env *gasket_api.CreateEnvRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", EnvTable)

	row := r.db.QueryRow(query, env.Title)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *GasketRepo) UpdateEnv(ctx context.Context, env *gasket_api.UpdateEnvRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", EnvTable)

	_, err := r.db.Exec(query, env.Title, env.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) DeleteEnv(ctx context.Context, env *gasket_api.DeleteEnvRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", EnvTable)

	if _, err := r.db.Exec(query, env.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

//---
func (r *GasketRepo) GetEnvData(ctx context.Context, gasketId string) (env []models.EnvDataDTO, err error) {
	query := fmt.Sprintf(`SELECT id, env_id, gasket_id, m, specific_pres FROM %s WHERE gasket_id=$1 ORDER BY env_id`, EnvDataTable)

	if err := r.db.Select(&env, query, gasketId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return env, nil
}

func (r *GasketRepo) CreateManyEnvData(ctx context.Context, data *gasket_api.CreateManyEnvDataRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (env_id, gasket_id, m, specific_pres) VALUES ($1, $2, $3, $4)", EnvDataTable)

	args := make([]interface{}, 0)
	args = append(args, data.Data[0].EnvId, data.GasketId, data.Data[0].M, data.Data[0].SpecificPres)

	for i, d := range data.Data {
		if i > 0 {
			query += fmt.Sprintf(", ($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
			args = append(args, d.EnvId, data.GasketId, d.M, d.SpecificPres)
		}
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) CreateEnvData(ctx context.Context, data *gasket_api.CreateEnvDataRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (env_id, gasket_id, m, specific_pres) VALUES ($1, $2, $3, $4)", EnvDataTable)

	if _, err := r.db.Exec(query, data.EnvId, data.GasketId, data.M, data.SpecificPres); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) UpdateEnvData(ctx context.Context, data *gasket_api.UpdateEnvDataRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if data.EnvId != "" {
		setValues = append(setValues, fmt.Sprintf("env_id=$%d", argId))
		args = append(args, data.EnvId)
		argId++
	}
	if data.GasketId != "" {
		setValues = append(setValues, fmt.Sprintf("gasket_id=$%d", argId))
		args = append(args, data.GasketId)
		argId++
	}

	if data.M != 0 {
		setValues = append(setValues, fmt.Sprintf("m=$%d", argId))
		args = append(args, data.M)
		argId++
	}
	if data.SpecificPres != 0 {
		setValues = append(setValues, fmt.Sprintf("specific_pres=$%d", argId))
		args = append(args, data.SpecificPres)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", EnvDataTable, setQuery, argId)

	args = append(args, data.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *GasketRepo) DeleteEnvData(ctx context.Context, data *gasket_api.DeleteEnvDataRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", EnvDataTable)

	if _, err := r.db.Exec(query, data.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
