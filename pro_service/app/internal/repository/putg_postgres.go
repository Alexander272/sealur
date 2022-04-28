package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type PutgRepo struct {
	db *sqlx.DB
}

func NewPutgRepo(db *sqlx.DB) *PutgRepo {
	return &PutgRepo{db: db}
}

func (r *PutgRepo) Get(req *proto.GetPutgRequest) (putg []models.Putg, err error) {
	query := fmt.Sprintf(`SELECT id, type_fl_id, type_pr, form, construction, temperatures, reinforce, obturator, i_limiter, o_limiter, 
		coating, mounting, graphite FROM %s WHERE form=$1 AND flange_id=$2`, PutgTable)

	if err = r.db.Select(&putg, query, req.Form, req.FlangeId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return putg, nil
}

func (r *PutgRepo) Create(putg models.PutgDTO) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (flange_id, type_fl_id, type_pr, form, construction, temperatures, reinforce, obturator, 
		i_limiter, o_limiter, coating, mounting, graphite) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) 
		RETURNING id`, PutgTable)

	flangeId, err := strconv.Atoi(putg.FlangeId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	typeFlId, err := strconv.Atoi(putg.TypeFlId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	row := r.db.QueryRow(query, flangeId, typeFlId, putg.TypePr, putg.Form, putg.Construction, putg.Temperatures, putg.Reinforce,
		putg.Obturator, putg.ILimiter, putg.OLimiter, putg.Coating, putg.Mounting, putg.Graphite)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *PutgRepo) Update(putg models.PutgDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET type_fl_id=$1, type_pr=$2, form=$3, construction=$4, temperatures=$5, reinforce=$6,	obturator=$7,
		i_limiter=$8, o_limiter=$9, coating=$10, mounting=$11, graphite=$12, flange_id=$13 WHERE id=$14`, PutgTable)

	id, err := strconv.Atoi(putg.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, putg.TypeFlId, putg.TypePr, putg.Form, putg.Construction, putg.Temperatures, putg.Reinforce,
		putg.Obturator, putg.ILimiter, putg.OLimiter, putg.Coating, putg.Mounting, putg.Graphite, putg.FlangeId, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PutgRepo) Delete(putg *proto.DeletePutgRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", PutgTable)

	id, err := strconv.Atoi(putg.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
