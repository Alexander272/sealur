package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type MatRepo struct {
	db *sqlx.DB
}

func NewMatRepo(db *sqlx.DB) *MatRepo {
	return &MatRepo{db: db}
}

func (r *MatRepo) GetAll(*pro_api.GetMaterialsRequest) (mats []models.Materials, err error) {
	query := fmt.Sprintf("SELECT id, title, type_mat FROM %s ORDER BY id", MaterialsTable)

	if err = r.db.Select(&mats, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return mats, nil
}

func (r *MatRepo) Create(mat *pro_api.CreateMaterialsRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title, type_mat) VALUES ($1, $2) RETURNING id", MaterialsTable)
	row := r.db.QueryRow(query, mat.Title, mat.TypeMat)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return fmt.Sprintf("%d", idInt), nil
}

func (r *MatRepo) Update(mat *pro_api.UpdateMaterialsRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1, type_mat=$2 WHERE id=$3", MaterialsTable)

	id, err := strconv.Atoi(mat.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, mat.Title, mat.TypeMat, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *MatRepo) Delete(mat *pro_api.DeleteMaterialsRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", MaterialsTable)

	id, err := strconv.Atoi(mat.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
