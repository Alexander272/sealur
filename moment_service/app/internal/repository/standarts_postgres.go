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
	query := fmt.Sprintf("SELECT id, title, type_id, title_dn, title_pn, is_need_row, rows FROM %s WHERE type_id=$1 ORDER BY id", StandartsTable)

	if err := r.db.Select(&standarts, query, req.TypeId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return standarts, nil
}

func (r *FlangeRepo) CreateStandart(ctx context.Context, stand *moment_api.CreateStandartRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (title, type_id, title_dn, title_pn, is_need_row, rows) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, StandartsTable)

	row := r.db.QueryRow(query, stand.Title, stand.TypeId, stand.TitleDn, stand.TitlePn, stand.IsNeedRow, strings.Join(stand.Rows, "; "))
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
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if stand.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, stand.Title)
		argId++
	}
	if stand.TypeId != "" {
		setValues = append(setValues, fmt.Sprintf("type_id=$%d", argId))
		args = append(args, stand.TypeId)
		argId++
	}
	if stand.TitleDn != "" {
		setValues = append(setValues, fmt.Sprintf("title_dn=$%d", argId))
		args = append(args, stand.TitleDn)
		argId++
	}
	if stand.TitlePn != "" {
		setValues = append(setValues, fmt.Sprintf("title_pn=$%d", argId))
		args = append(args, stand.TitlePn)
		argId++
	}
	if stand.IsNeedRow {
		setValues = append(setValues, fmt.Sprintf("is_need_row=$%d", argId))
		args = append(args, stand.IsNeedRow)
		argId++
	}
	if stand.Rows != nil {
		setValues = append(setValues, fmt.Sprintf("rows=$%d", argId))
		args = append(args, strings.Join(stand.Rows, "; "))
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", StandartsTable, setQuery, argId)

	args = append(args, stand.Id)
	_, err := r.db.Exec(query, args...)
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
