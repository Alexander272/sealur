package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
	"github.com/jmoiron/sqlx"
)

type GasketRepo struct {
	db *sqlx.DB
}

func NewGasketRepo(db *sqlx.DB) *GasketRepo {
	return &GasketRepo{db: db}
}

func (r *GasketRepo) GetFullData(ctx context.Context, req models.GetGasket) (data models.FullDataGasket, err error) {
	query := fmt.Sprintf(`SELECT %s.id, %s.title as gasket_title, %s.title as env_title, permissible_pres, compression, epsilon, thickness, m, specific_pres, 
		%s.title as type_title, %s.label as type_label
		FROM %s
		INNER JOIN %s ON %s.type_id = %s.id
		INNER JOIN %s ON %s.gasket_id = %s.gasket_id
		INNER JOIN %s ON %s.gasket_id = %s.id
		INNER JOIN %s ON %s.env_id = %s.id
		WHERE %s.gasket_id=$1 AND env_id=$2 AND thickness=$3`,
		GasketDataTable, GasketTable, EnvTable,
		TypeGasketTable, TypeGasketTable,
		GasketDataTable,
		TypeGasketTable, GasketDataTable, TypeGasketTable,
		EnvDataTable, GasketDataTable, EnvDataTable,
		GasketTable, GasketDataTable, GasketTable,
		EnvTable, EnvDataTable, EnvTable, GasketDataTable,
	)

	if err := r.db.Get(&data, query, req.GasketId, req.EnvId, req.Thickness); err != nil {
		return data, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *GasketRepo) GetGasket(ctx context.Context, req *moment_api.GetGasketRequest) (gasket []models.GasketDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title FROM %s ORDER BY id`, GasketTable)

	if err := r.db.Select(&gasket, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return gasket, nil
}

func (r *GasketRepo) GetGasketWithThick(ctx context.Context, req *moment_api.GetGasketRequest) (gasket []models.GasketWithThick, err error) {
	query := fmt.Sprintf(`SELECT %s.id, title, thickness FROM %s
		INNER JOIN %s ON %s.id = %s.gasket_id ORDER BY id`,
		GasketTable, GasketTable, GasketDataTable, GasketTable, GasketDataTable)

	if err := r.db.Select(&gasket, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return gasket, nil
}

func (r *GasketRepo) CreateGasket(ctx context.Context, gasket *moment_api.CreateGasketRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", GasketTable)

	row := r.db.QueryRow(query, gasket.Title)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *GasketRepo) UpdateGasket(ctx context.Context, gasket *moment_api.UpdateGasketRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", GasketTable)

	_, err := r.db.Exec(query, gasket.Title, gasket.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) DeleteGasket(ctx context.Context, gasket *moment_api.DeleteGasketRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", GasketTable)

	if _, err := r.db.Exec(query, gasket.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

//---

func (r *GasketRepo) CreateGasketData(ctx context.Context, data *moment_api.CreateGasketDataRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (gasket_id, permissible_pres, compression, epsilon, thickness, type_id)
		VALUES ($1, $2, $3, $4, $5, $6)`, GasketDataTable)

	if _, err := r.db.Exec(query, data.GasketId, data.PermissiblePres, data.Compression, data.Epsilon, data.Thickness, data.TypeId); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) UpdateGasketData(ctx context.Context, data *moment_api.UpdateGasketDataRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if data.GasketId != "" {
		setValues = append(setValues, fmt.Sprintf("gasket_id=$%d", argId))
		args = append(args, data.GasketId)
		argId++
	}

	if data.PermissiblePres != 0 {
		setValues = append(setValues, fmt.Sprintf("permissible_pres=$%d", argId))
		args = append(args, data.PermissiblePres)
		argId++
	}
	if data.Compression != 0 {
		setValues = append(setValues, fmt.Sprintf("compression=$%d", argId))
		args = append(args, data.Compression)
		argId++
	}
	if data.Epsilon != 0 {
		setValues = append(setValues, fmt.Sprintf("epsilon=$%d", argId))
		args = append(args, data.Epsilon)
		argId++
	}
	if data.Thickness != 0 {
		setValues = append(setValues, fmt.Sprintf("thickness=$%d", argId))
		args = append(args, data.Thickness)
		argId++
	}
	if data.TypeId != "" {
		setValues = append(setValues, fmt.Sprintf("type_id=$%d", argId))
		args = append(args, data.TypeId)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", GasketDataTable, setQuery, argId)

	args = append(args, data.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *GasketRepo) DeleteGasketData(ctx context.Context, data *moment_api.DeleteGasketDataRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", GasketDataTable)

	if _, err := r.db.Exec(query, data.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
