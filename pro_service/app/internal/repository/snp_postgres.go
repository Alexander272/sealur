package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
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
	query := fmt.Sprintf(`SELECT id, type_fl_id, type_pr, filler, materials, def_mat, mounting, graphite 
		FROM %s WHERE stand_id=$1 AND flange_id=$2`, SNPTable)

	var data []models.SNP
	if err = r.db.Select(&data, query, req.StandId, req.FlangeId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, d := range data {
		s := proto.SNP(d)
		snp = append(snp, &s)
	}

	return snp, nil
}

func (r *SNPRepo) Create(snp *proto.CreateSNPRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (stand_id, flange_id, type_fl_id, type_pr, filler, materials, def_mat, mounting, graphite) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, SNPTable)

	standId, err := strconv.Atoi(snp.StandId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	row := r.db.QueryRow(query, standId, snp.FlangeId, snp.TypeFlId, snp.TypePr, snp.Fillers, snp.Materials, snp.DefMat,
		snp.Mounting, snp.Graphite)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *SNPRepo) Update(snp *proto.UpdateSNPRequest) error {
	query := fmt.Sprintf(`UPDATE %s SET stand_id=$1, type_fl_id=$2, type_pr=$3, filler=$4, materials=$5, def_mat=$6,
		mounting=$7, graphite=$8, flange_id=$9 WHERE id=$10`, SNPTable)

	id, err := strconv.Atoi(snp.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	standId, err := strconv.Atoi(snp.StandId)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, standId, snp.TypeFlId, snp.TypePr, snp.Fillers, snp.Materials, snp.DefMat,
		snp.Mounting, snp.Graphite, snp.FlangeId, id)
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
