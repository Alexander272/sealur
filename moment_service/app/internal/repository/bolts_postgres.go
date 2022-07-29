package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (r *FlangeRepo) GetBolts(ctx context.Context, req *moment_api.GetBoltsRequest) (bolts []models.BoltsDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title, diameter, area FROM %s ORDER BY diameter`, BoltsTable)

	if err := r.db.Select(&bolts, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return bolts, nil
}

func (r *FlangeRepo) CreateBolt(ctx context.Context, bolt *moment_api.CreateBoltRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (title, diameter, area) VALUES ($1, $2, $3)", BoltsTable)

	if _, err := r.db.Exec(query, bolt.Title, bolt.Diameter, bolt.Area); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) UpdateBolt(ctx context.Context, bolt *moment_api.UpdateBoltRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if bolt.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, bolt.Title)
		argId++
	}

	if bolt.Diameter != 0 {
		setValues = append(setValues, fmt.Sprintf("diameter=$%d", argId))
		args = append(args, bolt.Diameter)
		argId++
	}
	if bolt.Area != 0 {
		setValues = append(setValues, fmt.Sprintf("area=$%d", argId))
		args = append(args, bolt.Area)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", BoltsTable, setQuery, argId)

	args = append(args, bolt.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *FlangeRepo) DeleteBolt(ctx context.Context, bolt *moment_api.DeleteBoltRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", BoltsTable)

	if _, err := r.db.Exec(query, bolt.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
