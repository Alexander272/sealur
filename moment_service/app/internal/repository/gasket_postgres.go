package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/jmoiron/sqlx"
)

type GasketRepo struct {
	db *sqlx.DB
}

func NewGasketRepo(db *sqlx.DB) *GasketRepo {
	return &GasketRepo{db: db}
}

func (r *GasketRepo) Get(ctx context.Context, gasket models.GetGasket) (g models.Gasket, err error) {
	// query := fmt.Sprintf("SELECT ")

	return g, nil
}
