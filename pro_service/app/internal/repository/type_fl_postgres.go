package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type TypeFlRepo struct {
	db *sqlx.DB
}

func NewTypeFlRepo(db *sqlx.DB) *TypeFlRepo {
	return &TypeFlRepo{db: db}
}

func (r *TypeFlRepo) Get() (fl []*pro_api.TypeFl, err error) {
	query := fmt.Sprintf("SELECT id, title, descr, short FROM %s WHERE basis=true", TypeFLTable)

	if err = r.db.Select(&fl, query); err != nil {
		return fl, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return fl, nil
}

func (r *TypeFlRepo) GetAll() (fl []*pro_api.TypeFl, err error) {
	query := fmt.Sprintf("SELECT id, title, descr, short FROM %s", TypeFLTable)

	if err = r.db.Select(&fl, query); err != nil {
		return fl, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return fl, nil
}

func (r *TypeFlRepo) Create(fl *pro_api.CreateTypeFlRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title, descr, short, basis) VALUES ($1, $2, $3, $4) RETURNING id", TypeFLTable)
	row := r.db.QueryRow(query, fl.Title, fl.Descr, fl.Short, fl.Basis)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *TypeFlRepo) Update(fl *pro_api.UpdateTypeFlRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1, descr=$2, short=$3, basis=$4 WHERE id=$5", TypeFLTable)

	id, err := strconv.Atoi(fl.Id)
	if err != nil {
		return fmt.Errorf("failed to covert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, fl.Title, fl.Descr, fl.Short, fl.Basis, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *TypeFlRepo) Delete(fl *pro_api.DeleteTypeFlRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TypeFLTable)

	id, err := strconv.Atoi(fl.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
