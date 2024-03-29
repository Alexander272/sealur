package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type SNPRepo struct {
	db *sqlx.DB
}

func NewSNPRepo(db *sqlx.DB) *SNPRepo {
	return &SNPRepo{db: db}
}

func (r *SNPRepo) Get(req *pro_api.GetSNPRequest) (snp []models.SNP, err error) {
	query := fmt.Sprintf(`SELECT id, type_fl_id, type_pr, filler, frame, in_ring, ou_ring, mounting, graphite 
		FROM %s WHERE stand_id=$1 AND flange_id=$2 ORDER BY type_pr DESC`, SNPTable)

	if err = r.db.Select(&snp, query, req.StandId, req.FlangeId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return snp, nil
}

func (r *SNPRepo) Create(snp models.SnpDTO) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (stand_id, flange_id, type_fl_id, type_pr, filler, frame, in_ring, ou_ring, mounting, graphite) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, SNPTable)

	standId, err := strconv.Atoi(snp.StandId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	flangeId, err := strconv.Atoi(snp.FlangeId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	typeFlId, err := strconv.Atoi(snp.TypeFlId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	row := r.db.QueryRow(query, standId, flangeId, typeFlId, snp.TypePr, snp.Fillers, snp.Frame, snp.Ir, snp.Or,
		snp.Mounting, snp.Graphite)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *SNPRepo) Update(snp models.SnpDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET stand_id=$1, type_fl_id=$2, type_pr=$3, filler=$4, frame=$5, in_ring=$6, ou_ring=$7,
		mounting=$8, graphite=$9, flange_id=$10 WHERE id=$11`, SNPTable)

	id, err := strconv.Atoi(snp.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	standId, err := strconv.Atoi(snp.StandId)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, standId, snp.TypeFlId, snp.TypePr, snp.Fillers, snp.Frame, snp.Ir, snp.Or,
		snp.Mounting, snp.Graphite, snp.FlangeId, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SNPRepo) Delete(snp *pro_api.DeleteSNPRequest) error {
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

func (r *SNPRepo) GetByCondition(cond string) (snp []models.SNP, err error) {
	query := fmt.Sprintf(`SELECT id, filler, frame, in_ring, ou_ring, mounting, graphite FROM %s WHERE %s`, SNPTable, cond)

	if err = r.db.Select(&snp, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return snp, nil
}

func (r *SNPRepo) UpdateAddit(snp models.UpdateAdditDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET filler=$1, frame=$2, in_ring=$3, ou_ring=$4, mounting=$5, graphite=$6 WHERE id=$7`, SNPTable)

	id, err := strconv.Atoi(snp.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, snp.Fillers, snp.Frame, snp.Ir, snp.Or, snp.Mounting, snp.Graphite, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
