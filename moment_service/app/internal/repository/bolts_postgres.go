package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (r *FlangeRepo) GetBolts(ctx context.Context, req *moment_api.GetBoltsRequest) (bolts []models.BoltsDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title, diameter, area, is_inch FROM %s WHERE is_inch=$1 ORDER BY diameter`, BoltsTable)

	if err := r.db.Select(&bolts, query, req.IsInch); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return bolts, nil
}

func (r *FlangeRepo) GetAllBolts(ctx context.Context, req *moment_api.GetBoltsRequest) (bolts []models.BoltsDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title, diameter, area, is_inch FROM %s ORDER BY diameter`, BoltsTable)

	if err := r.db.Select(&bolts, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return bolts, nil
}

func (r *FlangeRepo) CreateBolt(ctx context.Context, bolt *moment_api.CreateBoltRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (title, diameter, area, is_inch) VALUES ($1, $2, $3, $4)", BoltsTable)

	if _, err := r.db.Exec(query, bolt.Title, bolt.Diameter, bolt.Area, bolt.IsInch); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) CreateBolts(ctx context.Context, bolts *moment_api.CreateBoltsRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	for i, b := range bolts.Bolts {
		setValues = append(setValues, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		args = append(args, b.Title, b.Diameter, b.Area, b.IsInch)
	}

	query := fmt.Sprintf("INSERT INTO %s (title, diameter, area, is_inch) VALUES %s", BoltsTable, strings.Join(setValues, ", "))

	if _, err := r.db.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) UpdateBolt(ctx context.Context, bolt *moment_api.UpdateBoltRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1, diameter=$2, area=$3, is_inch=$4 WHERE id=$5", BoltsTable)

	_, err := r.db.Exec(query, bolt.Title, bolt.Diameter, bolt.Area, bolt.IsInch, bolt.Id)
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
