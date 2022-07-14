package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/jmoiron/sqlx"
)

type Flange interface {
	GetSize(ctx context.Context, dn, pn float64, standId string) (models.FlangeSize, error)
}

type Materials interface {
	GetAll(ctx context.Context, markId string) ([]models.Materials, error)
}

type Gasket interface {
	Get(ctx context.Context, gasket models.GetGasket) (models.Gasket, error)
}

type Repositories struct {
	Flange
	Materials
	Gasket
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{
		Flange:    NewFlangeRepo(db),
		Materials: NewMaterialsRepo(db),
		Gasket:    NewGasketRepo(db),
	}
}
