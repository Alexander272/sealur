package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type FlangeRepo struct {
	db *sqlx.DB
}

func NewFlangeRepo(db *sqlx.DB) *FlangeRepo {
	return &FlangeRepo{db: db}
}

func (r *FlangeRepo) GetAll() (flanges []*proto.Flange, err error) {
	query := fmt.Sprintf("SELECT id, title, short FROM %s", FlangeTable)

	if err = r.db.Select(&flanges, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return flanges, nil
}

func (r *FlangeRepo) GetByTitle(title, short string) (flange *proto.Flange, err error) {
	query := fmt.Sprintf("SELECT id, title, short from %s WHERE lower(title)=lower($1) OR lower(short)=lower($2)", FlangeTable)

	if err := r.db.Get(&flange, query, title, short); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return flange, nil
}

// TODO дописать создание таблицы с размерами
func (r *FlangeRepo) Create(fl *proto.CreateFlangeRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title, short) VALUES ($1, $2) RETURNING id", FlangeTable)
	row := r.db.QueryRow(query, fl.Title, fl.Short)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *FlangeRepo) Update(fl *proto.UpdateFlangeRequest) error {
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

// TODO дописать удаление таблицы с размерами
func (r *FlangeRepo) Delete(fl *proto.DeleteFlangeRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", FlangeTable)

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
