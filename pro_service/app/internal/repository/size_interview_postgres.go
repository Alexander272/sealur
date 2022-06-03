package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type SizeIntRepo struct {
	db *sqlx.DB
}

func NewSizeIntRepo(db *sqlx.DB) *SizeIntRepo {
	return &SizeIntRepo{db: db}
}

func (r *SizeIntRepo) Get(req *proto.GetSizesIntRequest) (sizes []models.SizeInterview, err error) {
	// TODO дописать join и сделать запрос
	// var query string

	// query = fmt.Sprintf(`SELECT id, dy, py, d1, d2, d, h1, h2, bolt, count_bolt FROM %s
	// 	WHERE stand_id=$1 AND type_fl_id=$2 ORDER BY count`, req.FlangeId, req.TypeFl)

	// var data []models.Size
	// if err = r.db.Select(&data, query, req.StandId, req.TypeFlId); err != nil {
	// 	return nil, fmt.Errorf("failed to execute query. error: %w", err)
	// }

	return sizes, nil
}

func (r *SizeIntRepo) Create(size *proto.CreateSizeIntRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (count, type_fl_id, flange_id, dy, py, d1, d2, d, h1, h2, bolt, count_bolt) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`, SizeIntrTable)

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

	row := r.db.QueryRow(query, size.Number, size.TypeFl, size.FlangeId, size.Dy, size.Py, size.D1, size.D2, size.D, size.H1,
		size.H2, size.BoltId, size.CountBolt)

	var idInt int
	if err = row.Scan(&idInt); err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *SizeIntRepo) Update(size *proto.UpdateSizeIntRequest) error {
	query := fmt.Sprintf(`UPDATE %s SET dy=$1, py=$2, flange_id=$3, type_fl_id=$4, d1=$5, d2=$6, d=$7, h1=$8, h2=$9, bolt_id=$10,
		count_bolt=$11 WHERE id=$12`, SizeIntrTable)

	id, err := strconv.Atoi(size.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	// standId, err := strconv.Atoi(size.StandId)
	// if err != nil {
	// 	return fmt.Errorf("failed to convert string to int. error: %w", err)
	// }

	_, err = r.db.Exec(query, size.Dy, size.Py, size.FlangeId, size.TypeFl, size.D1, size.D2, size.D, size.H1, size.H2, size.BoltId,
		size.CountBolt, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SizeIntRepo) Delete(size *proto.DeleteSizeIntRequest) error {
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

// func (r *SizeIntRepo) DeleteAll(size *proto.DeleteAllSizeIntRequest) error {
// 	query := fmt.Sprintf("DELETE FROM %s WHERE ", SizeIntrTable)

// 	_, err := r.db.Exec(query)
// 	if err != nil {
// 		return fmt.Errorf("failed to execute query. error: %w", err)
// 	}
// 	return nil
// }
