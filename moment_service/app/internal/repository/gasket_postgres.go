package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type GasketRepo struct {
	db *sqlx.DB
}

func NewGasketRepo(db *sqlx.DB) *GasketRepo {
	return &GasketRepo{db: db}
}

func (r *GasketRepo) GetFullData(ctx context.Context, req models.GetGasket) (data models.FullDataGasket, err error) {
	// query := fmt.Sprintf("SELECT ")
	//TODO
	return data, nil
}

func (r *GasketRepo) GetGasket(ctx context.Context, req *moment_proto.GetGasketRequest) (gasket []models.GasketDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title FROM %s ORDER BY id`, GasketTable)

	if err := r.db.Select(&gasket, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return gasket, nil
}

func (r *GasketRepo) CreateGasket(ctx context.Context, gasket *moment_proto.CreateGasketRequest) (id string, err error) {
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

func (r *GasketRepo) UpdateGasket(ctx context.Context, gasket *moment_proto.UpdateGasketRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", GasketTable)

	_, err := r.db.Exec(query, gasket.Title, gasket.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) DeleteGasket(ctx context.Context, gasket *moment_proto.DeleteGasketRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", GasketTable)

	if _, err := r.db.Exec(query, gasket.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

//---

func (r *GasketRepo) GetTypeGasket(ctx context.Context, req *moment_proto.GetGasketTypeRequest) (types []models.TypeGasketDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title FROM %s ORDER BY id`, TypeGasketTable)

	if err := r.db.Select(&types, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return types, nil
}

func (r *GasketRepo) CreateTypeGasket(ctx context.Context, typeGasket *moment_proto.CreateGasketTypeRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", TypeGasketTable)

	row := r.db.QueryRow(query, typeGasket.Title)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *GasketRepo) UpdateTypeGasket(ctx context.Context, typeGasket *moment_proto.UpdateGasketTypeRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", TypeGasketTable)

	_, err := r.db.Exec(query, typeGasket.Title, typeGasket.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) DeleteTypeGasket(ctx context.Context, typeGasket *moment_proto.DeleteGasketTypeRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TypeGasketTable)

	if _, err := r.db.Exec(query, typeGasket.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

//---

func (r *GasketRepo) GetEnv(ctx context.Context, req *moment_proto.GetEnvRequest) (env []models.EnvDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title FROM %s ORDER BY id`, EnvTable)

	if err := r.db.Select(&env, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return env, nil
}

func (r *GasketRepo) CreateEnv(ctx context.Context, env *moment_proto.CreateEnvRequest) (id string, err error) {
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

func (r *GasketRepo) UpdateEnv(ctx context.Context, env *moment_proto.UpdateEnvRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", EnvTable)

	_, err := r.db.Exec(query, env.Title, env.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) DeleteEnv(ctx context.Context, env *moment_proto.DeleteEnvRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", EnvTable)

	if _, err := r.db.Exec(query, env.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

//---

func (r *GasketRepo) CreateEnvData(ctx context.Context, data *moment_proto.CreateEnvDataRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (env_id, gasket_id, m, specific_pres) VALUES ($1, $2, $3, $4)", EnvDataTable)

	if _, err := r.db.Exec(query, data.EnvId, data.GasketId, data.M, data.SpecificPres); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) UpdateEnvData(ctx context.Context, data *moment_proto.UpdateEnvDataRequest) error {
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

func (r *GasketRepo) DeleteEnvData(ctx context.Context, data *moment_proto.DeleteEnvDataRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", EnvDataTable)

	if _, err := r.db.Exec(query, data.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

//---

func (r *GasketRepo) CreateGasketData(ctx context.Context, data *moment_proto.CreateGasketDataRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (gasket_id, permissible_pres, compression, epsilon, thickness, type_id)
		VALUES ($1, $2, $3, $4, $5, $6)`, GasketDataTable)

	if _, err := r.db.Exec(query, data.GasketId, data.PermissiblePres, data.Compression, data.Epsilon, data.Thickness, data.TypeId); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *GasketRepo) UpdateGasketData(ctx context.Context, data *moment_proto.UpdateGasketDataRequest) error {
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

func (r *GasketRepo) DeleteGasketData(ctx context.Context, data *moment_proto.DeleteGasketDataRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", GasketDataTable)

	if _, err := r.db.Exec(query, data.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
