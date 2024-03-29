package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type SizesRepo struct {
	db *sqlx.DB
}

func NewSizesRepo(db *sqlx.DB) *SizesRepo {
	return &SizesRepo{db: db}
}

func (r *SizesRepo) Get(req *pro_api.GetSizesRequest) (sizes []*pro_api.Size, err error) {
	var query string
	if strings.Contains(strings.ToLower(req.TypePr), "путг") {
		query = fmt.Sprintf(`SELECT id, dn, pn, d4, d3, d2, d1, h, s2, s3, type_pr, type_fl_id FROM size_%s WHERE LOWER(type_pr) = LOWER('%s') 
		AND (stand_id=$1 OR stand_id=0) AND type_fl_id=$2 ORDER BY count`, req.Flange, req.TypePr)
	} else {
		query = fmt.Sprintf(`SELECT id, dn, pn, d4, d3, d2, d1, h, s2, s3, type_pr, type_fl_id FROM size_%s WHERE LOWER(type_pr) LIKE LOWER('%%%s%%') 
		AND (stand_id=$1 OR stand_id=0) AND type_fl_id=$2 ORDER BY count`, req.Flange, req.TypePr)
	}

	var data []models.Size
	if err = r.db.Select(&data, query, req.StandId, req.TypeFlId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, d := range data {
		sizes = append(sizes, &pro_api.Size{
			Id:       d.Id,
			Dn:       d.Dn,
			Pn:       d.Pn,
			D4:       d.D4,
			D3:       d.D3,
			D2:       d.D2,
			D1:       d.D1,
			H:        d.H,
			S2:       d.S2,
			S3:       d.S3,
			TypePr:   d.TypePr,
			TypeFlId: d.TypeFlId,
		})
	}

	return sizes, nil
}

func (r *SizesRepo) GetAll(req *pro_api.GetSizesRequest) (sizes []*pro_api.Size, err error) {
	var query string
	if strings.Contains(strings.ToLower(req.TypePr), "путг") {
		query = fmt.Sprintf(`SELECT id, dn, pn, d4, d3, d2, d1, h, s2, s3, type_pr, type_fl_id FROM size_%s WHERE LOWER(type_pr) like LOWER('%s%%') 
		AND (stand_id=$1 OR stand_id=0) AND type_fl_id=$2 ORDER BY type_pr, count`, req.Flange, req.TypePr)
	}

	var data []models.Size
	if err = r.db.Select(&data, query, req.StandId, req.TypeFlId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, d := range data {
		sizes = append(sizes, &pro_api.Size{
			Id:       d.Id,
			Dn:       d.Dn,
			Pn:       d.Pn,
			D4:       d.D4,
			D3:       d.D3,
			D2:       d.D2,
			D1:       d.D1,
			H:        d.H,
			S2:       d.S2,
			S3:       d.S3,
			TypePr:   d.TypePr,
			TypeFlId: d.TypeFlId,
		})
	}

	return sizes, nil
}

func (r *SizesRepo) Create(size *pro_api.CreateSizeRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO size_%s (count, dn, pn, type_fl_id, type_pr, stand_id, d4, d3, d2, d1, h, s2, s3) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`, size.Flange)

	standId, err := strconv.Atoi(size.StandId)
	if err != nil {
		return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	if size.Number == 0 {
		var max int32
		query := fmt.Sprintf(`SELECT MAX(count) FROM size_%s WHERE LOWER(type_pr) LIKE LOWER('%%%s%%') 
		AND (stand_id=$1 OR stand_id=0) AND type_fl_id=$2`, size.Flange, size.TypePr)

		if err := r.db.Get(&max, query, size.StandId, size.TypeFlId); err != nil {
			return "", fmt.Errorf("failed to get max count query. error: %w", err)
		}
		size.Number = max + 1
	}

	row := r.db.QueryRow(query, size.Number, size.Dn, size.Pn, size.TypeFlId, size.TypePr, standId, size.D4, size.D3, size.D2, size.D1,
		size.H, size.S2, size.S3)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *SizesRepo) Update(size *pro_api.UpdateSizeRequest) error {
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

	_, err = r.db.Exec(query, size.Dn, size.Pn, size.TypePr, standId, size.D4, size.D3, size.D2, size.D1, size.H, size.TypeFlId,
		size.S2, size.S3, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SizesRepo) Delete(size *pro_api.DeleteSizeRequest) error {
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

func (r *SizesRepo) DeleteAll(size *pro_api.DeleteAllSizeRequest) error {
	query := fmt.Sprintf("DELETE FROM size_%s WHERE LOWER(type_pr) LIKE LOWER('%%%s%%')", size.Flange, size.TypePr)

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
