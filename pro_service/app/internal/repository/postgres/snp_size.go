package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type SnpSizeRepo struct {
	db *sqlx.DB
}

func NewSnpSizeRepo(db *sqlx.DB) *SnpSizeRepo {
	return &SnpSizeRepo{db: db}
}

func (r *SnpSizeRepo) Get(ctx context.Context)
