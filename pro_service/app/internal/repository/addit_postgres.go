package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type AdditRepo struct {
	db *sqlx.DB
}

func NewAdditRepo(db *sqlx.DB) *AdditRepo {
	return &AdditRepo{db: db}
}

func (r *AdditRepo) GetAll() (addit []*proto.Additional, err error) {
	query := fmt.Sprintf("SELECT id, materials, mod, temperature, mounting, graphite FROM %s LIMIT 1", AdditionalTable)

	if err = r.db.Select(&addit, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return addit, nil
}

func (r *AdditRepo) Create(add *proto.CreateAddRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (materials, mod, temperature, mounting, graphite)
		VALUES ($1, $2, $3, $4, $5)`, AdditionalTable)

	_, err := r.db.Exec(query, add.Materials, add.Mod, add.Temperature, add.Mounting, add.Graphite)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateMat(mat *proto.UpdateAddMatRequest) error {
	query := fmt.Sprintf("UPDATE %s SET materials=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(mat.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, mat.Materials, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateMod(mod *proto.UpdateAddModRequest) error {
	query := fmt.Sprintf("UPDATE %s SET mod=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(mod.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, mod.Mod, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateTemp(temp *proto.UpdateAddTemRequest) error {
	query := fmt.Sprintf("UPDATE %s SET temperature=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(temp.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, temp.Temperature, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateMoun(moun *proto.UpdateAddMounRequest) error {
	query := fmt.Sprintf("UPDATE %s SET mounting=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(moun.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error %w", err)
	}

	_, err = r.db.Exec(query, moun.Mounting, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AdditRepo) UpdateGrap(grap *proto.UpdateAddGrapRequest) error {
	query := fmt.Sprintf("UPDATE %s SET graphite=$1 WHERE id=$2", AdditionalTable)

	id, err := strconv.Atoi(grap.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, grap.Graphite, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
