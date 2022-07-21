package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type FlangeRepo struct {
	db *sqlx.DB
}

func NewFlangeRepo(db *sqlx.DB) *FlangeRepo {
	return &FlangeRepo{db: db}
}

func (r *FlangeRepo) GetFlangeSize(ctx context.Context, req *moment_proto.GetFlangeSizeRequest) (size models.FlangeSize, err error) {
	//TODO
	return size, nil
}

func (r *FlangeRepo) CreateFlangeSize(ctx context.Context, size *moment_proto.CreateFlangeSizeRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (stand_id, pn, d, d6, d_out, h, s0, s1, lenght, count, bolt_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, FlangeSizeTable)

	_, err := r.db.Exec(query, size.StandId, size.Pn, size.D, size.D6, size.DOut, size.H, size.S0, size.S1, size.Lenght, size.Count, size.BoltId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) UpdateFlangeSize(ctx context.Context, size *moment_proto.UpdateFlangeSizeRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if size.StandId != "" {
		setValues = append(setValues, fmt.Sprintf("stand_id=$%d", argId))
		args = append(args, size.StandId)
		argId++
	}

	if size.Pn != 0 {
		setValues = append(setValues, fmt.Sprintf("pn=$%d", argId))
		args = append(args, size.Pn)
		argId++
	}
	if size.D != 0 {
		setValues = append(setValues, fmt.Sprintf("d=$%d", argId))
		args = append(args, size.D)
		argId++
	}
	if size.D6 != 0 {
		setValues = append(setValues, fmt.Sprintf("d6=$%d", argId))
		args = append(args, size.D6)
		argId++
	}
	if size.DOut != 0 {
		setValues = append(setValues, fmt.Sprintf("d_out=$%d", argId))
		args = append(args, size.DOut)
		argId++
	}
	if size.H != 0 {
		setValues = append(setValues, fmt.Sprintf("h=$%d", argId))
		args = append(args, size.H)
		argId++
	}
	if size.S0 != 0 {
		setValues = append(setValues, fmt.Sprintf("s0=$%d", argId))
		args = append(args, size.S0)
		argId++
	}
	if size.S1 != 0 {
		setValues = append(setValues, fmt.Sprintf("s1=$%d", argId))
		args = append(args, size.S1)
		argId++
	}
	if size.Lenght != 0 {
		setValues = append(setValues, fmt.Sprintf("lenght=$%d", argId))
		args = append(args, size.Lenght)
		argId++
	}
	if size.Count != 0 {
		setValues = append(setValues, fmt.Sprintf("count=$%d", argId))
		args = append(args, size.Count)
		argId++
	}
	if size.BoltId != "" {
		setValues = append(setValues, fmt.Sprintf("bolt_id=$%d", argId))
		args = append(args, size.BoltId)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", FlangeSizeTable, setQuery, argId)

	args = append(args, size.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *FlangeRepo) DeleteFlangeSize(ctx context.Context, size *moment_proto.DeleteFlangeSizeRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", FlangeSizeTable)

	if _, err := r.db.Exec(query, size.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
