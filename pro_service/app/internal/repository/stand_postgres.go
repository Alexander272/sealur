package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type StandRepo struct {
	db *sqlx.DB
}

func NewStandRepo(db *sqlx.DB) *StandRepo {
	return &StandRepo{db: db}
}

func (r *StandRepo) GetAll(stand proto.GetStands) (stands []proto.Stand, err error) {
	query := fmt.Sprintf("SELECT id, title FROM %s", StandTable)

	if err = r.db.Select(&stands, query); err != nil {
		logger.Error("failed to execute query. error: %w", err)
		return stands, err
	}
	return stands, nil
}

func (r *StandRepo) Create(stand proto.CreateStand) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", StandTable)
	res, err := r.db.Exec(query, stand.Title)
	if err != nil {
		logger.Error("failed to execute query. error: %w", err)
		return id, err
	}

	idInt, err := res.LastInsertId()
	if err != nil {
		logger.Error("failed to get id. error: %w", err)
		return id, err
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *StandRepo) Update(stand proto.UpdateStand) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", StandTable)

	id, err := strconv.Atoi(stand.Id)
	if err != nil {
		logger.Error("falied to convert string to int. error: %w", err)
		return err
	}

	_, err = r.db.Exec(query, stand.Title, id)
	if err != nil {
		logger.Error("failed to execute query. error: %w", err)
		return err
	}

	return nil
}

func (r *StandRepo) Delete(stand proto.DeleteStand) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", StandTable)

	id, err := strconv.Atoi(stand.Id)
	if err != nil {
		logger.Error("failed to convert string to int. error: %w", err)
		return err
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		logger.Error("failed to execute query. error: %w", err)
		return err
	}

	return nil
}
