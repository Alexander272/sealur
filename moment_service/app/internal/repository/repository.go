package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/jmoiron/sqlx"
)

type Flange interface {
	GetSize(ctx context.Context, dn, pn float32) (models.FlangeSize, error)
}

type Materials interface {
	Get(ctx context.Context, markId, tempId string) (models.Materials, error)
}

type Repositories struct {
	Flange
	Materials
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{
		Flange:    NewFlangeRepo(db),
		Materials: NewMaterialsRepo(db),
	}
}
