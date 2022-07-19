package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type MaterialsRepo struct {
	db *sqlx.DB
}

func NewMaterialsRepo(db *sqlx.DB) *MaterialsRepo {
	return &MaterialsRepo{db: db}
}

func (r *MaterialsRepo) GetMaterials(ctx context.Context, req *moment_proto.GetMaterialsRequest) (materials []models.MaterialsDTO, err error) {
	query := fmt.Sprintf("SELECT id, title FROM %s ORDER BY id", MaterialsTable)

	if err := r.db.Select(&materials, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return materials, nil
}

func (r *MaterialsRepo) GetAllData(ctx context.Context, markId string) (materials models.MaterialsAll, err error) {

	return materials, nil
}

func (r *MaterialsRepo) CreateMaterial(ctx context.Context, material *moment_proto.CreateMaterialRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", MaterialsTable)

	row := r.db.QueryRow(query, material.Title)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *MaterialsRepo) UpdateMaterial(ctx context.Context, material *moment_proto.UpdateMaterialRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", MaterialsTable)

	_, err := r.db.Exec(query, material.Title, material.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialsRepo) DeleteMaterial(ctx context.Context, material *moment_proto.DeleteMaterialRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", MaterialsTable)

	if _, err := r.db.Exec(query, material.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialsRepo) CreateVoltage(ctx context.Context, voltage *moment_proto.CreateVoltageRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (mark_id, temperature, voltage) VALUES ($1, $2, $3)", VoltageTable)

	// row := r.db.QueryRow(query, voltage.MarkId, voltage.Temperature, voltage.Voltage)
	// if row.Err() != nil {
	// 	return "", fmt.Errorf("failed to execute query. error: %w", err)
	// }

	// var idInt int
	// if err := row.Scan(&idInt); err != nil {
	// 	return "", fmt.Errorf("failed to scan result. error: %w", err)
	// }

	// return fmt.Sprintf("%d", idInt), nil

	args := make([]interface{}, 0)
	args = append(args, voltage.MarkId, voltage.Voltage[0].Temperature, voltage.Voltage[0].Voltage)

	for i, v := range voltage.Voltage {
		if i > 0 {
			query += fmt.Sprintf(", ($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
			args = append(args, voltage.MarkId, v.Temperature, v.Voltage)
		}
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialsRepo) UpdateVoltage(ctx context.Context, voltage *moment_proto.UpdateVoltageRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if voltage.MarkId != "" {
		setValues = append(setValues, fmt.Sprintf("mark_id=$%d", argId))
		args = append(args, voltage.MarkId)
		argId++
	}
	if voltage.Temperature != 0 {
		setValues = append(setValues, fmt.Sprintf("temperature=$%d", argId))
		args = append(args, voltage.Temperature)
		argId++
	}
	if voltage.Voltage != 0 {
		setValues = append(setValues, fmt.Sprintf("voltage=$%d", argId))
		args = append(args, voltage.Voltage)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", VoltageTable, setQuery, argId)

	args = append(args, voltage.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *MaterialsRepo) DeleteVoltage(ctx context.Context, voltage *moment_proto.DeleteVoltageRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", VoltageTable)

	if _, err := r.db.Exec(query, voltage.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MaterialsRepo) CreateElasticity(ctx context.Context, elasticity *moment_proto.CreateElasticityRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (mark_id, temperature, elasticity) VALUES ($1, $2, $3) RETURNING id", ElasticityTable)

	// row := r.db.QueryRow(query, elasticity.MarkId, elasticity.Temperature, elasticity.Elasticity)
	// if row.Err() != nil {
	// 	return "", fmt.Errorf("failed to execute query. error: %w", err)
	// }

	// var idInt int
	// if err := row.Scan(&idInt); err != nil {
	// 	return "", fmt.Errorf("failed to scan result. error: %w", err)
	// }

	// return fmt.Sprintf("%d", idInt), nil

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

func (r *MaterialsRepo) CreateAlpha(ctx context.Context, alpha *moment_proto.CreateAlphaRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (mark_id, temperature, alpha) VALUES ($1, $2, $3) RETURNING id", AlphaTable)

	// row := r.db.QueryRow(query, alpha.MarkId, alpha.Temperature, alpha.Alpha)
	// if row.Err() != nil {
	// 	return "", fmt.Errorf("failed to execute query. error: %w", err)
	// }

	// var idInt int
	// if err := row.Scan(&idInt); err != nil {
	// 	return "", fmt.Errorf("failed to scan result. error: %w", err)
	// }

	// return fmt.Sprintf("%d", idInt), nil

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

func (r *MaterialsRepo) UpateAlpha(ctx context.Context, alpha *moment_proto.UpdateAlphaRequest) error {
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

func (r *MaterialsRepo) DeleteAlpha(ctx context.Context, alpha *moment_proto.DeleteAlphaRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", AlphaTable)

	if _, err := r.db.Exec(query, alpha.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
