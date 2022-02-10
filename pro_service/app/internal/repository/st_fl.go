package repository

import (
	"fmt"
	"strconv"

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
	//TODO дописать join и изменить proto
	query := fmt.Sprintf("SELECT id, stand_id, fl_ids FROM %s", StFLTable)

	if err := r.db.Select(&st, query); err != nil {
		return st, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return st, nil
}

func (r *StFlRepo) Create(st *proto.CreateStFlRequest) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s (stand_id, fl_ids) VALUES ($1, $2)", StFLTable)

	row := r.db.QueryRow(query, st.StandId, st.FlangeId)

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *StFlRepo) Update(st *proto.UpdateStFlRequest) error {
	query := fmt.Sprintf("UPDATE %s SET stand_id=$1, fl_ids=$2 WHERE id=$3", StFLTable)

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
