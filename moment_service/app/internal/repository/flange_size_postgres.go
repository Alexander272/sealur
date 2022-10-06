package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/moment/flange_api"
	"github.com/jmoiron/sqlx"
)

type FlangeRepo struct {
	db *sqlx.DB
}

func NewFlangeRepo(db *sqlx.DB) *FlangeRepo {
	return &FlangeRepo{db: db}
}

func (r *FlangeRepo) GetFlangeSize(ctx context.Context, req *flange_api.GetFlangeSizeRequest) (size models.FlangeSize, err error) {
	query := fmt.Sprintf(`SELECT %s.id, pn, d, d6, d_out, x, a, h, s0, s1, length, count, diameter, area FROM %s
		INNER JOIN %s on bolt_id=%s.id WHERE stand_id=$1 AND dn=$2 AND pn=$3 AND row=$4`,
		FlangeSizeTable, FlangeSizeTable, BoltsTable, BoltsTable)

	if err := r.db.Get(&size, query, req.StandId, req.Dn, req.Pn, req.Row); err != nil {
		return size, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return size, nil
}

func (r *FlangeRepo) GetBasisFlangeSizes(ctx context.Context, req models.GetBasisSize) (sizes []models.FlangeSize, err error) {
	args := make([]interface{}, 0)
	args = append(args, req.StandId)

	orderField := "d"
	if req.IsInch {
		orderField = "dmm"
	}

	var query string
	if req.IsUseRow {
		query = fmt.Sprintf(`SELECT id, pn, d, dn, COALESCE(d = 0, true) as is_empty_d FROM %s 
			WHERE stand_id=$1 AND row=$2 ORDER BY %s, pn`, FlangeSizeTable, orderField)
		args = append(args, req.Row)
	} else {
		query = fmt.Sprintf(`SELECT id, pn, d, dn, COALESCE(d = 0, true) as is_empty_d FROM %s 
			WHERE stand_id=$1 ORDER BY %s, pn`, FlangeSizeTable, orderField)
	}

	if err := r.db.Select(&sizes, query, args...); err != nil {
		return sizes, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return sizes, nil
}

func (r *FlangeRepo) GetFullFlangeSize(ctx context.Context, req *flange_api.GetFullFlangeSizeRequest, row int32) (size []models.FlangeSizeDTO, err error) {
	query := fmt.Sprintf(`SELECT id, stand_id, pn, dn, dmm, d, d6, d_out, x, a, h, s0, s1, length, count, bolt_id 
		FROM %s WHERE stand_id=$1 AND row=$2`, FlangeSizeTable)

	if err := r.db.Select(&size, query, req.StandId, row); err != nil {
		return size, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return size, nil
}

func (r *FlangeRepo) CreateFlangeSize(ctx context.Context, size *flange_api.CreateFlangeSizeRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (stand_id, pn, dn, dmm, d, d6, d_out, x, a, h, s0, s1, length, count, bolt_id, row)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`, FlangeSizeTable)

	_, err := r.db.Exec(query, size.StandId, size.Pn, size.Dn, size.Dmm, size.D, size.D6, size.DOut, size.X, size.A, size.H, size.S0, size.S1,
		size.Length, size.Count, size.BoltId, size.Row)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) CreateFlangeSizes(ctx context.Context, size *flange_api.CreateFlangeSizesRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	for i, s := range size.Sizes {
		logger.Debug(s)
		setValues = append(setValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*16+1, i*16+2, i*16+3, i*16+4, i*16+5, i*16+6, i*16+7, i*16+8, i*16+9, i*16+10, i*16+11, i*16+12, i*16+13, i*16+14, i*16+15, i*16+16))
		args = append(args, s.StandId, s.Pn, s.Dn, s.Dmm, s.D, s.D6, s.DOut, s.X, s.A, s.H, s.S0, s.S1, s.Length, s.Count, s.BoltId, s.Row)
	}

	query := fmt.Sprintf("INSERT INTO %s (stand_id, pn, dn, dmm, d, d6, d_out, x, a, h, s0, s1, length, count, bolt_id, row) VALUES %s",
		FlangeSizeTable, strings.Join(setValues, ", "))

	if _, err := r.db.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) UpdateFlangeSize(ctx context.Context, size *flange_api.UpdateFlangeSizeRequest) error {
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
	if size.Dn != "" {
		setValues = append(setValues, fmt.Sprintf("dn=$%d", argId))
		args = append(args, size.Dn)
		argId++
	}
	if size.Dmm != 0 {
		setValues = append(setValues, fmt.Sprintf("dmm=$%d", argId))
		args = append(args, size.Dmm)
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
	if size.X != 0 {
		setValues = append(setValues, fmt.Sprintf("x=$%d", argId))
		args = append(args, size.X)
		argId++
	}
	if size.A != 0 {
		setValues = append(setValues, fmt.Sprintf("a=$%d", argId))
		args = append(args, size.A)
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
	if size.Length != 0 {
		setValues = append(setValues, fmt.Sprintf("length=$%d", argId))
		args = append(args, size.Length)
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

func (r *FlangeRepo) DeleteFlangeSize(ctx context.Context, size *flange_api.DeleteFlangeSizeRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", FlangeSizeTable)

	if _, err := r.db.Exec(query, size.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
