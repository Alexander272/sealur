package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
)

func (r *FlangeRepo) GetTypeFlange(ctx context.Context, req *moment_proto.GetTypeFlangeRequest) (typeFlange []models.TypeFlangeDTO, err error) {
	query := fmt.Sprintf(`SELECT id, title FROM %s ORDER BY title`, TypeFlangeTable)

	if err := r.db.Select(&typeFlange, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return typeFlange, nil
}

func (r *FlangeRepo) CreateTypeFlange(ctx context.Context, typeFlange *moment_proto.CreateTypeFlangeRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", TypeFlangeTable)

	row := r.db.QueryRow(query, typeFlange.Title)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *FlangeRepo) UpdateTypeFlange(ctx context.Context, typeFlange *moment_proto.UpdateTypeFlangeRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", TypeFlangeTable)

	_, err := r.db.Exec(query, typeFlange.Title, typeFlange.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) DeleteTypeFlange(ctx context.Context, typeFlange *moment_proto.DeleteTypeFlangeRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", TypeFlangeTable)

	if _, err := r.db.Exec(query, typeFlange.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FlangeRepo) GetStandarts(ctx context.Context, req *moment_proto.GetStandartsRequest) (standarts []models.StandartDTO, err error) {
	query := fmt.Sprintf("SELECT id, title, type_id FROM %s WHERE type_id=$1 ORDER BY id", StandartsTable)

	if err := r.db.Select(&standarts, query, req.TypeId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return standarts, nil
}

func (r *FlangeRepo) CreateStandart(ctx context.Context, stand *moment_proto.CreateStandartRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title, type_id) VALUES ($1, $2) RETURNING id", TypeFlangeTable)

	row := r.db.QueryRow(query, stand.Title, stand.TypeId)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var idInt int
	if err := row.Scan(&idInt); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", idInt), nil
}

func (r *FlangeRepo) UpdateStandart(ctx context.Context, stand *moment_proto.UpdateStandartRequest) error {
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
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", StandartsTable, setQuery, argId)

	args = append(args, stand.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *FlangeRepo) DeleteStandart(ctx context.Context, stand *moment_proto.DeleteStandartRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", StandartsTable)

	if _, err := r.db.Exec(query, stand.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
