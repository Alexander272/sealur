package repository

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/jmoiron/sqlx"
)

type FlangeRepo struct {
	db *sqlx.DB
}

func NewFlangeRepo(db *sqlx.DB) *FlangeRepo {
	return &FlangeRepo{db: db}
}

func (r *FlangeRepo) GetSize(ctx context.Context, dn, pn float64, standId string) (size models.FlangeSize, err error) {

	return size, nil
}
