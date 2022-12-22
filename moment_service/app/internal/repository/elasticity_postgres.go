package repository

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur_proto/api/moment/material_api"
)

func (r *MaterialsRepo) CreateElasticity(ctx context.Context, elasticity *material_api.CreateElasticityRequest) error {
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

func (r *MaterialsRepo) UpdateElasticity(ctx context.Context, elasticity *material_api.UpdateElasticityRequest) error {
	query := fmt.Sprintf("UPDATE %s SET temperature=$1, elasticity=$2 WHERE id=$3", ElasticityTable)

	_, err := r.db.Exec(query, elasticity.Temperature, elasticity.Elasticity, elasticity.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *MaterialsRepo) DeleteElasticity(ctx context.Context, elasticity *material_api.DeleteElasticityRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", ElasticityTable)

	if _, err := r.db.Exec(query, elasticity.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
