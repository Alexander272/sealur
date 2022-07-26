package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type StandRepo struct {
	db *sqlx.DB
}

func NewStandRepo(db *sqlx.DB) *StandRepo {
	return &StandRepo{db: db}
}

func (r *StandRepo) GetAll(stand *pro_api.GetStandsRequest) (stands []*pro_api.Stand, err error) {
	query := fmt.Sprintf("SELECT id, title FROM %s", StandTable)

	if err = r.db.Select(&stands, query); err != nil {
		return stands, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return stands, nil
}

func (r *StandRepo) GetByTitle(title string) (stand []*pro_api.Stand, err error) {
	query := fmt.Sprintf("SELECT id, title FROM %s WHERE lower(title)=lower($1) limit 1", StandTable)

	if err = r.db.Select(&stand, query, title); err != nil {
		return stand, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return stand, nil
}

func (r *StandRepo) Create(stand *pro_api.CreateStandRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", StandTable)
	row := r.db.QueryRow(query, stand.Title)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *StandRepo) Update(stand *pro_api.UpdateStandRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", StandTable)

	id, err := strconv.Atoi(stand.Id)
	if err != nil {
		return fmt.Errorf("falied to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, stand.Title, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *StandRepo) Delete(stand *pro_api.DeleteStandRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", StandTable)

	id, err := strconv.Atoi(stand.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
