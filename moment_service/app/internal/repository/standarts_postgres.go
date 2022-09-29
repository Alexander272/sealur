package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment_api"
)

func (r *FlangeRepo) GetTypeFlange(ctx context.Context, req *moment_api.GetTypeFlangeRequest) (typeFlange []models.TypeFlangeDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title, label FROM %s ORDER BY id`, TypeFlangeTable)

	if err := r.db.Select(&typeFlange, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return typeFlange, nil
}

func (r *FlangeRepo) CreateTypeFlange(ctx context.Context, typeFlange *moment_api.CreateTypeFlangeRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title, label) VALUES ($1, $2) RETURNING id", TypeFlangeTable)

	row := r.db.QueryRow(query, typeFlange.Title, typeFlange.Label)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *FlangeRepo) UpdateTypeFlange(ctx context.Context, typeFlange *moment_api.UpdateTypeFlangeRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if typeFlange.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, typeFlange.Title)
		argId++
	}
	if typeFlange.Label != "" {
		setValues = append(setValues, fmt.Sprintf("label=$%d", argId))
		args = append(args, typeFlange.Label)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", TypeFlangeTable, setQuery, argId)

	args = append(args, typeFlange.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *FlangeRepo) DeleteTypeFlange(ctx context.Context, typeFlange *moment_api.DeleteTypeFlangeRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TypeFlangeTable)

	if _, err := r.db.Exec(query, typeFlange.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) GetStandarts(ctx context.Context, req *moment_api.GetStandartsRequest) (standarts []models.StandartDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title, type_id, title_dn, title_pn, is_need_row, rows, is_inch, has_designation 
		FROM %s WHERE type_id=$1 ORDER BY id`, StandartsTable)

	if err := r.db.Select(&standarts, query, req.TypeId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return standarts, nil
}

func (r *FlangeRepo) CreateStandart(ctx context.Context, stand *moment_api.CreateStandartRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (title, type_id, title_dn, title_pn, is_need_row, rows, is_inch, has_designation) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, StandartsTable)

	row := r.db.QueryRow(query, stand.Title, stand.TypeId, stand.TitleDn, stand.TitlePn, stand.IsNeedRow,
		strings.Join(stand.Rows, "; "), stand.IsInch, stand.HasDesignation)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", row.Err())
	}

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *FlangeRepo) UpdateStandart(ctx context.Context, stand *moment_api.UpdateStandartRequest) error {
	query := fmt.Sprintf(`UPDATE %s SET title=$1, type_id=$2, title_dn=$3, title_pn=$4, is_need_row=$5, rows=$6, 
		is_inch=$7, has_desigantion=$8 WHERE id=$9`, StandartsTable)

	_, err := r.db.Exec(query, stand.Title, stand.TypeId, stand.TitleDn, stand.TitlePn, stand.IsNeedRow,
		strings.Join(stand.Rows, "; "), stand.IsInch, stand.HasDesignation, stand.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *FlangeRepo) DeleteStandart(ctx context.Context, stand *moment_api.DeleteStandartRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", StandartsTable)

	if _, err := r.db.Exec(query, stand.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
