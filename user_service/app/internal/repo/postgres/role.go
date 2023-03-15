package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/user/models/role_model"
	"github.com/jmoiron/sqlx"
)

type RoleRepo struct {
	db *sqlx.DB
}

func NewRoleRepo(db *sqlx.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) GetDefault(ctx context.Context) (*role_model.Role, error) {
	var data models.Role
	query := fmt.Sprintf(`SELECT id, title, code FROM %s WHERE is_default=true LIMIT 1`, RoleTable)

	if err := r.db.Get(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	role := &role_model.Role{
		Id:    data.Id,
		Title: data.Title,
		Code:  data.Code,
	}
	return role, nil
}
