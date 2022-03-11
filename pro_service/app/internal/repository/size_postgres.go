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
	query := fmt.Sprintf(`SELECT id, dn, pn, d4, d3, d2, d1, h, s2, s3 FROM size_%s WHERE LOWER(type_pr) LIKE LOWER('%%%s%%') 
		AND stand_id=$1 AND type_fl_id=$2 ORDER BY dn`, req.Flange, req.TypePr)

	if err = r.db.Select(&sizes, query, req.StandId, req.TypeFlId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return sizes, nil
}

func (r *SizesRepo) Create(size *proto.CreateSizeRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO size_%s (dn, pn, type_fl_id, type_pr, stand_id, d4, d3, d2, d1, h, s2, s3) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`, size.Flange)

	standId, err := strconv.Atoi(size.StandId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	row := r.db.QueryRow(query, size.Dn, size.Pn, size.TypeFlId, size.TypePr, standId, size.D4, size.D3, size.D2, size.D1, size.H, size.S2, size.S3)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *SizesRepo) Update(size *proto.UpdateSizeRequest) error {
	query := fmt.Sprintf(`UPDATE size_%s SET dn=$1, pn=$2, type_pr=$3, stand_id=$4, d4=$5, d3=$6, d2=$7, d1=$8, h=$9, type_fl_id=$10,
		s2=$11, s3=$12 WHERE id=$13`, size.Flange)

	id, err := strconv.Atoi(size.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	standId, err := strconv.Atoi(size.StandId)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, size.Dn, size.Pn, size.TypePr, standId, size.D4, size.D3, size.D2, size.D1, size.H, size.TypeFlId, size.S2, size.S3, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SizesRepo) Delete(size *proto.DeleteSizeRequest) error {
	query := fmt.Sprintf("DELETE FROM size_%s WHERE id=$1", size.Flange)

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
