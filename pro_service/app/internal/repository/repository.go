package repository

import "github.com/jmoiron/sqlx"

type Repositories struct {
}

func NewRepo(db *sqlx.DB) *Repositories {
	return &Repositories{}
}
