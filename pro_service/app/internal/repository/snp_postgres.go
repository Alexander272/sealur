package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type SNPRepo struct {
	db *sqlx.DB
}

func NewSNPRepo(db *sqlx.DB) *SNPRepo {
	return &SNPRepo{db: db}
}

func (r *SNPRepo) Get(req *proto.GetSNPRequest) (snp []*proto.SNP, err error) {
	query := fmt.Sprintf("SELECT id, type_p as typeP, fillers, materials FROM %s WHERE stand_id=$1 AND type_fl=$2", SNPTable)

	if err = r.db.Select(&snp, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return snp, nil
}

func (r *SNPRepo) Create(snp *proto.CreateSNPRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (stand_id, type_fl, type_p, fillers, materials) VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`, SNPTable)

	standId, err := strconv.Atoi(snp.StandId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	res, err := r.db.Exec(query, standId, snp.TypeFl, snp.TypeP, snp.Fillers, snp.Materials)
	if err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	idInt, err := res.LastInsertId()
	if err != nil {
		return id, fmt.Errorf("failed to get id. error: %w", err)
	}
	return fmt.Sprintf("%d", idInt), nil
}

func (r *SNPRepo) Update(snp *proto.UpdateSNPRequest) error {
	query := fmt.Sprintf("UPDATE %s SET stand_id=$1, type_fl=$2, type_p=$3, fillers=$4, materials=$5 WHERE id=$6", SNPTable)

	id, err := strconv.Atoi(snp.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	standId, err := strconv.Atoi(snp.StandId)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, standId, snp.TypeFl, snp.TypeP, snp.Fillers, snp.Materials, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPRepo) Delete(snp *proto.DeleteSNPRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", SNPTable)

	id, err := strconv.Atoi(snp.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
