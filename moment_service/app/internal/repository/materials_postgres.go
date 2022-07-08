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

func (r *MaterialsRepo) GetAll(ctx context.Context, markId string) ([]models.Materials, error) {
	materials := []models.Materials{
		{MarkId: "1", Temp: 20, Voltage: 130.0, Elasticity: 2.13, Alpha: 0},
		{MarkId: "1", Temp: 100, Voltage: 126.0, Elasticity: 2.10, Alpha: 11.1},
	}

	return materials, nil
}
