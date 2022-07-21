package repository

import (
	"context"
	"fmt"

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
	query := fmt.Sprintf(`SELECT id, title FROM %s WHERE 
			(SELECT count(mark_id) FROM %s GROUP BY mark_id HAVING mark_id = %s.id) > 0 AND
			(SELECT count(mark_id) FROM %s GROUP BY mark_id HAVING mark_id = %s.id) > 0 AND
			(SELECT count(mark_id) FROM %s GROUP BY mark_id HAVING mark_id = %s.id) > 0
		ORDER BY id`, MaterialsTable, ElasticityTable, MaterialsTable, VoltageTable, MaterialsTable, AlphaTable, MaterialsTable)

	if err := r.db.Select(&materials, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return materials, nil
}

func (r *MaterialsRepo) GetMaterialsWithIsEmpty(ctx context.Context, req *moment_proto.GetMaterialsRequest,
) (materials []models.MaterialsWithIsEmpty, err error) {
	query := fmt.Sprintf(`SELECT id, title, 
			(SELECT count(mark_id) FROM %s GROUP BY mark_id HAVING mark_id = %s.id) = 0 as is_empty_elasticity, 
			(SELECT count(mark_id) FROM %s GROUP BY mark_id HAVING mark_id = %s.id) = 0 as is_empty_voltage, 
			(SELECT count(mark_id) FROM %s GROUP BY mark_id HAVING mark_id = %s.id) = 0 as is_empty_alpha
		FROM %s ORDER BY id`, ElasticityTable, MaterialsTable, VoltageTable, MaterialsTable, AlphaTable, MaterialsTable, MaterialsTable)

	if err := r.db.Select(&materials, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return materials, nil
}

func (r *MaterialsRepo) GetAllData(ctx context.Context, req *moment_proto.GetMaterialsDataRequest) (materials models.MaterialsAll, err error) {
	voltageQuery := fmt.Sprintf("SELECT id, temperature, voltage FROM %s WHERE mark_id=$1 ORDER BY temperature", VoltageTable)
	var voltage []models.Voltage

	if err := r.db.Select(&voltage, voltageQuery, req.MarkId); err != nil {
		return materials, fmt.Errorf("failed to execute query. error: %w", err)
	}

	elasticityQuery := fmt.Sprintf("SELECT id, temperature, elasticity FROM %s WHERE mark_id=$1 ORDER BY temperature", ElasticityTable)
	var elasticity []models.Elasticity

	if err := r.db.Select(&elasticity, elasticityQuery, req.MarkId); err != nil {
		return materials, fmt.Errorf("failed to execute query. error: %w", err)
	}

	alphaQuery := fmt.Sprintf("SELECT id, temperature, alpha FROM %s WHERE mark_id=$1 ORDER BY temperature", AlphaTable)
	var alpha []models.Alpha

	if err := r.db.Select(&alpha, alphaQuery, req.MarkId); err != nil {
		return materials, fmt.Errorf("failed to execute query. error: %w", err)
	}

	materials.Voltage = voltage
	materials.Elasticity = elasticity
	materials.Alpha = alpha

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
