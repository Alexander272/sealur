package repository

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type StFlRepo struct {
	db *sqlx.DB
}

func NewStFlRepo(db *sqlx.DB) *StFlRepo {
	return &StFlRepo{db: db}
}

func (r *StFlRepo) Get() (st []*proto.StFl, err error) {
	query := fmt.Sprintf(`SELECT st_fl.id, stand_id, stand.title AS stand, fl_id, flange.title AS flange, short FROM %s 
		LEFT JOIN %s ON (stand_id = stand.id) LEFT JOIN %s ON (fl_id = flange.id)`, StFLTable, StandTable, FlangeTable)

	var data []models.StFl
	if err := r.db.Select(&data, query); err != nil {
		return st, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, d := range data {
		s := proto.StFl(d)
		st = append(st, &s)
	}

	return st, nil
}

func (r *StFlRepo) Create(st *proto.CreateStFlRequest) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s (stand_id, fl_id) VALUES ($1, $2) RETURNING id", StFLTable)

	standId, err := strconv.Atoi(st.StandId)
	if err != nil {
		return "", fmt.Errorf("failed to convert string to int. error: %w", err)
	}
	flangeId, err := strconv.Atoi(st.FlangeId)
	if err != nil {
		return "", fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	row := r.db.QueryRow(query, standId, flangeId)

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *StFlRepo) Update(st *proto.UpdateStFlRequest) error {
	query := fmt.Sprintf("UPDATE %s SET stand_id=$1, fl_id=$2 WHERE id=$3", StFLTable)

	id, err := strconv.Atoi(st.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, st.StandId, st.FlangeId, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *StFlRepo) Delete(st *proto.DeleteStFlRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", StFLTable)

	id, err := strconv.Atoi(st.Id)
	if err != nil {
		return fmt.Errorf("failed to convert string to int. error: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
