package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type PutgmRepo struct {
	db *sqlx.DB
}

func NewPutgmRepo(db *sqlx.DB) *PutgmRepo {
	return &PutgmRepo{db: db}
}

func (r *PutgmRepo) Get(req *proto.GetPutgmRequest) (putgm []models.Putgm, err error) {
	query := fmt.Sprintf(`SELECT id, type_fl_id, type_pr, form, construction, temperatures, basis, obturator, coating, mounting, 
	graphite FROM %s WHERE form=$1 AND flange_id=$2 ORDER BY type_fl_id`, PutgmTable)

	if err = r.db.Select(&putgm, query, req.Form, req.FlangeId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return putgm, nil
}

func (r *PutgmRepo) Create(putgm models.PutgmDTO) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (flange_id, type_fl_id, type_pr, form, construction, temperatures, basis, obturator, 
		coating, mounting, graphite) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING id`, PutgmTable)

	flangeId, err := strconv.Atoi(putgm.FlangeId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	typeFlId, err := strconv.Atoi(putgm.TypeFlId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	row := r.db.QueryRow(query, flangeId, typeFlId, putgm.TypePr, putgm.Form, putgm.Construction, putgm.Temperatures, putgm.Basis,
		putgm.Obturator, putgm.Coating, putgm.Mounting, putgm.Graphite)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *PutgmRepo) Update(putgm models.PutgmDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET type_fl_id=$1, type_pr=$2, form=$3, construction=$4, temperatures=$5, basis=$6,	obturator=$7,
		coating=$8, mounting=$9, graphite=$10, flange_id=$11 WHERE id=$12`, PutgmTable)

	id, err := strconv.Atoi(putgm.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, putgm.TypeFlId, putgm.TypePr, putgm.Form, putgm.Construction, putgm.Temperatures, putgm.Basis,
		putgm.Obturator, putgm.Coating, putgm.Mounting, putgm.Graphite, putgm.FlangeId, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PutgmRepo) Delete(putgm *proto.DeletePutgmRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", PutgmTable)

	id, err := strconv.Atoi(putgm.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PutgmRepo) GetByCondition(cond string) (putgm []models.Putgm, err error) {
	query := fmt.Sprintf(`SELECT id, construction, temperatures, basis, obturator, coating, mounting, graphite FROM %s WHERE %s`,
		PutgmTable, cond)

	if err = r.db.Select(&putgm, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return putgm, nil
}

func (r *PutgmRepo) UpdateAddit(putgm models.UpdateAdditDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET construction=$1, temperatures=$2, basis=$3, obturator=$4, 
		coating=$5, mounting=$6, graphite=$7 WHERE id=$8`, PutgmTable)

	id, err := strconv.Atoi(putgm.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, putgm.Construction, putgm.Temperature, putgm.Basis, putgm.PObturator,
		putgm.Coating, putgm.Mounting, putgm.Graphite, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
