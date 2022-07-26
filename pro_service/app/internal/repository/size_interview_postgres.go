package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/jmoiron/sqlx"
)

type SizeIntRepo struct {
	db *sqlx.DB
}

func NewSizeIntRepo(db *sqlx.DB) *SizeIntRepo {
	return &SizeIntRepo{db: db}
}

func (r *SizeIntRepo) Get(req *pro_api.GetSizesIntRequest) (sizes []models.SizeInterview, err error) {
	query := fmt.Sprintf(`SELECT id, dy, py, d_up, d1, d2, d, h1, h2, bolt, count_bolt FROM %s
		WHERE flange_id=$1 AND type_fl_id=$2 AND row_count=$3 ORDER BY count`, SizeIntrTable)

	flangeId, err := strconv.Atoi(req.FlangeId)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	typeFl, err := strconv.Atoi(req.TypeFl)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	if err = r.db.Select(&sizes, query, flangeId, typeFl, req.Row); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return sizes, nil
}

func (r *SizeIntRepo) GetAll(req *pro_api.GetAllSizeIntRequest) (sizes []models.SizeInterview, err error) {
	query := fmt.Sprintf(`SELECT id, dy, py, d_up, d1, d2, d, h1, h2, bolt, count_bolt, row_count FROM %s
		WHERE flange_id=$1 AND type_fl_id=$2 ORDER BY count`, SizeIntrTable)

	flangeId, err := strconv.Atoi(req.FlangeId)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	typeFl, err := strconv.Atoi(req.TypeFl)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	if err = r.db.Select(&sizes, query, flangeId, typeFl); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return sizes, nil
}

func (r *SizeIntRepo) Create(size *pro_api.CreateSizeIntRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (count, type_fl_id, flange_id, dy, py, d_up, d1, d2, d, h1, h2, bolt, count_bolt, row_count) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`, SizeIntrTable)

	// flangeId, err := strconv.Atoi(size.FlangeId)
	// if err != nil {
	// 	return id, fmt.Errorf("failed to convert string to int. error: %w", err)
	// }

	if size.Number == 0 {
		var max int32
		query := fmt.Sprintf(`SELECT MAX(count) FROM %s WHERE flange_id=$1 AND type_fl_id=$2`, SizeIntrTable)

		if err := r.db.Get(&max, query, size.FlangeId, size.TypeFl); err != nil {
			return "", fmt.Errorf("failed to get max count query. error: %w", err)
		}
		size.Number = max + 1
	}

	row := r.db.QueryRow(query, size.Number, size.TypeFl, size.FlangeId, size.Dy, size.Py, size.DUp, size.D1, size.D2, size.D, size.H1,
		size.H2, size.Bolt, size.CountBolt, size.Row)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *SizeIntRepo) Update(size *pro_api.UpdateSizeIntRequest) error {
	query := fmt.Sprintf(`UPDATE %s SET dy=$1, py=$2, flange_id=$3, type_fl_id=$4, d_up=$5, d1=$6, d2=$7, d=$8, h1=$9, h2=$10, bolt=$11,
		count_bolt=$12, row_count=$13 WHERE id=$14`, SizeIntrTable)

	id, err := strconv.Atoi(size.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	// standId, err := strconv.Atoi(size.StandId)
	// if err != nil {
	// 	return fmt.Errorf("failed to convert string to int. error: %w", err)
	// }

	_, err = r.db.Exec(query, size.Dy, size.Py, size.FlangeId, size.TypeFl, size.DUp, size.D1, size.D2, size.D, size.H1, size.H2, size.Bolt,
		size.CountBolt, size.Row, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SizeIntRepo) Delete(size *pro_api.DeleteSizeIntRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", SizeIntrTable)

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

func (r *SizeIntRepo) DeleteAll(size *pro_api.DeleteAllSizeIntRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE flange_id=$1", SizeIntrTable)

	_, err := r.db.Exec(query, size.FlangeId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
