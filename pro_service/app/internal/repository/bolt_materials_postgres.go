package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type BoltMatRepo struct {
	db *sqlx.DB
}

func NewBoltMatRepo(db *sqlx.DB) *BoltMatRepo {
	return &BoltMatRepo{db: db}
}

func (r *BoltMatRepo) GetAll(*proto.GetBoltMaterialsRequest) (mats []models.BoltMaterials, err error) {
	query := fmt.Sprintf("SELECT id, title FROM %s", BoltsTable)

	if err = r.db.Select(&mats, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return mats, nil
}

func (r *BoltMatRepo) Create(mat *proto.CreateBoltMaterialsRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title, flange_id) VALUES ($1, $2) RETURNING id", BoltsTable)
	row := r.db.QueryRow(query, mat.Title, mat.FlangeId)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return fmt.Sprintf("%d", idInt), nil
}

func (r *BoltMatRepo) Update(mat *proto.UpdateBoltMaterialsRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1, flange_id=$2 WHERE id=$3", BoltsTable)

	id, err := strconv.Atoi(mat.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, mat.Title, mat.FlangeId, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *BoltMatRepo) Delete(mat *proto.DeleteBoltMaterialsRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", BoltsTable)

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
