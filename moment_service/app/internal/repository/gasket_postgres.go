package repository

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/jmoiron/sqlx"
)

type GasketRepo struct {
	db *sqlx.DB
}

func NewGasketRepo(db *sqlx.DB) *GasketRepo {
	return &GasketRepo{db: db}
}

var gaskets = []models.Gasket{
	{M: 3.0, SpecificPres: 69.0, PermissiblePres: 400.0, Compression: 1, Epsilon: 0.02 * float32(math.Pow10(5)), Thickness: 2.3},
	{M: 3.0, SpecificPres: 69.0, PermissiblePres: 400.0, Compression: 1, Epsilon: 0.02 * float32(math.Pow10(5)), Thickness: 3.2},
	{M: 3.0, SpecificPres: 69.0, PermissiblePres: 400.0, Compression: 1, Epsilon: 0.02 * float32(math.Pow10(5)), Thickness: 4.5},
	{M: 3.0, SpecificPres: 69.0, PermissiblePres: 400.0, Compression: 1, Epsilon: 0.02 * float32(math.Pow10(5)), Thickness: 6.5},
}

func (r *GasketRepo) Get(ctx context.Context, gasket models.GetGasket) (models.Gasket, error) {
	g := r.find(gasket.Thickness)
	return g, nil
}

func (r *GasketRepo) find(thic float32) models.Gasket {
	var res models.Gasket
	for _, g := range gaskets {
		if g.Thickness == thic {
			res = g
			break
		}
	}
	return res
}
