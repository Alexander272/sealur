package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type FlangeRepo struct {
	db *sqlx.DB
}

func NewFlangeRepo(db *sqlx.DB) *FlangeRepo {
	return &FlangeRepo{db: db}
}

func (r *FlangeRepo) GetAll() (flanges []*pro_api.Flange, err error) {
	query := fmt.Sprintf("SELECT id, title, short FROM %s", FlangeTable)

	if err = r.db.Select(&flanges, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return flanges, nil
}

func (r *FlangeRepo) GetByTitle(title, short string) (flange []*pro_api.Flange, err error) {
	query := fmt.Sprintf("SELECT id, title, short from %s WHERE lower(title)=lower($1) OR lower(short)=lower($2)", FlangeTable)

	if err := r.db.Select(&flange, query, title, short); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return flange, nil
}

func (r *FlangeRepo) Create(fl *pro_api.CreateFlangeRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title, short) VALUES ($1, $2) RETURNING id", FlangeTable)
	row := r.db.QueryRow(query, fl.Title, fl.Short)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	createTable := fmt.Sprintf(`CREATE TABLE size_%s (
		"id" serial not null unique,
		"count" int,
		"dn" int,
		"pn" text,
		"type_fl_id" int references type_fl (id) on delete cascade not null,
		"stand_id" int not null,
		"type_pr" text,
		"d4" text,
		"d3" text,
		"d2" text,
		"d1" text,
		"h" text,
		"s2" text,
		"s3" text
	);`, fl.Short)
	_, err = r.db.Exec(createTable)
	if err != nil {
		return id, fmt.Errorf("failed to execute query to create table. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *FlangeRepo) Update(fl *pro_api.UpdateFlangeRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1, short=$2 WHERE id=$3", FlangeTable)

	id, err := strconv.Atoi(fl.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, fl.Title, fl.Short, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) Delete(fl *pro_api.DeleteFlangeRequest) error {
	id, err := strconv.Atoi(fl.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	var tableName string

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING short", FlangeTable)
	row := r.db.QueryRow(query, id)

	if err = row.Scan(&tableName); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	dropTable := fmt.Sprintf("DROP TABLE size_%s", tableName)
	_, err = r.db.Exec(dropTable)
	if err != nil {
		return fmt.Errorf("failed to execute query to drop table. error: %w", err)
	}

	return nil
}
