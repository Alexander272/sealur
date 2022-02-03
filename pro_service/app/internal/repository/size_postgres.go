package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type SizesRepo struct {
	db *sqlx.DB
}

func NewSizesRepo(db *sqlx.DB) *SizesRepo {
	return &SizesRepo{db: db}
}

func (r *SizesRepo) Get(req *proto.GetSizesRequest) (sizes []*proto.Size, err error) {
	query := fmt.Sprintf("SELECT id, dn, pn, d4, d3, d2, d1, h FROM %s_%s WHERE type_p=$1 AND stand_id=$1", req.Flange, req.TypeFl)

	if err = r.db.Select(&sizes, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return sizes, nil
}

func (r *SizesRepo) Create(size *proto.CreateSizeRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s_%s (dn, pn, type_p, stand_id, d4, d3, d2, d1, h) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		size.Flange, size.TypeFl)

	standId, err := strconv.Atoi(size.StandId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	res, err := r.db.Exec(query, size.Dn, size.Pn, size.TypePr, standId, size.D4, size.D3, size.D2, size.D1, size.H)
	if err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	idInt, err := res.LastInsertId()
	if err != nil {
		return id, fmt.Errorf("failed to get id. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *SizesRepo) Update(size *proto.UpdateSizeRequest) error {
	query := fmt.Sprintf("UPDATE %s_%s SET dn=$1, pn=$2, type_p=$3, stand_id=$4, d4=$5, d3=$6, d2=$7, d1=$8, h=$9 WHERE id=$10",
		size.Flange, size.TypeFl)

	id, err := strconv.Atoi(size.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	standId, err := strconv.Atoi(size.StandId)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, size.Dn, size.Pn, size.TypePr, standId, size.D4, size.D3, size.D2, size.D1, size.H, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SizesRepo) Delete(size *proto.DeleteSizeRequest) error {
	query := fmt.Sprintf("DELETE FROM %s_%s WHERE id=$1", size.Flange, size.TypeFl)

	id, err := strconv.Atoi(size.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
