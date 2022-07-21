package repository

import (
	"context"
	"fmt"
	"strings"

	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (r *MaterialsRepo) CreateElasticity(ctx context.Context, elasticity *moment_proto.CreateElasticityRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (mark_id, temperature, elasticity) VALUES ($1, $2, $3)", ElasticityTable)

	args := make([]interface{}, 0)
	args = append(args, elasticity.MarkId, elasticity.Elasticity[0].Temperature, elasticity.Elasticity[0].Elasticity)

	for i, e := range elasticity.Elasticity {
		if i > 0 {
			query += fmt.Sprintf(", ($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
			args = append(args, elasticity.MarkId, e.Temperature, e.Elasticity)
		}
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialsRepo) UpdateElasticity(ctx context.Context, elasticity *moment_proto.UpdateElasticityRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if elasticity.MarkId != "" {
		setValues = append(setValues, fmt.Sprintf("mark_id=$%d", argId))
		args = append(args, elasticity.MarkId)
		argId++
	}
	if elasticity.Temperature != 0 {
		setValues = append(setValues, fmt.Sprintf("temperature=$%d", argId))
		args = append(args, elasticity.Temperature)
		argId++
	}
	if elasticity.Elasticity != 0 {
		setValues = append(setValues, fmt.Sprintf("elasticity=$%d", argId))
		args = append(args, elasticity.Elasticity)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", ElasticityTable, setQuery, argId)

	args = append(args, elasticity.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *MaterialsRepo) DeleteElasticity(ctx context.Context, elasticity *moment_proto.DeleteElasticityRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", ElasticityTable)

	if _, err := r.db.Exec(query, elasticity.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
