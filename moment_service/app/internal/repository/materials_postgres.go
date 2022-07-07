package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/jmoiron/sqlx"
)

type MaterialsRepo struct {
	db *sqlx.DB
}

func NewMaterialsRepo(db *sqlx.DB) *MaterialsRepo {
	return &MaterialsRepo{db: db}
}

func (r *MaterialsRepo) Get(ctx context.Context, markId, tempId string) (models.Materials, error) {
	material := models.Materials{}

	return material, nil
}
