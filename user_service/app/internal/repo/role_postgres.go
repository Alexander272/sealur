package repo

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/models"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type RoleRepo struct {
	db        *sqlx.DB
	tableName string
}

func NewRoleRepo(db *sqlx.DB, tableName string) *RoleRepo {
	return &RoleRepo{
		db:        db,
		tableName: tableName,
	}
}

func (r *RoleRepo) Get(ctx context.Context, req *proto_user.GetRolesRequest) (roles []models.Role, err error) {
	query := fmt.Sprintf("SELECT id, service, role FROM %s WHERE user_id=$1", r.tableName)

	if err := r.db.Select(&roles, query, req.UserId); err != nil {
		return roles, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return roles, nil
}

func (r *RoleRepo) GetAll(ctx context.Context, req *proto_user.GetAllRolesRequest) (roles []models.Role, err error) {
	query := fmt.Sprintf("SELECT id, user_id, service, role FROM %s ORDER BY user_id", r.tableName)

	if err := r.db.Select(&roles, query); err != nil {
		return roles, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return roles, nil
}

func (r *RoleRepo) Create(ctx context.Context, roles []*proto_user.CreateRoleRequest) error {
	// tx, err := r.db.Begin()
	// if err != nil {
	// 	return fmt.Errorf("failed to open transaction. error: %w", err)
	// }

	query := fmt.Sprintf("INSERT INTO %s (user_id, service, role) VALUES ($1, $2, $3)", r.tableName)
	args := make([]interface{}, 0)
	args = append(args, roles[0].UserId, roles[0].Service, roles[0].Role)

	for i, crr := range roles {
		if i > 0 {
			query += fmt.Sprintf(", ($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
			args = append(args, crr.UserId, crr.Service, crr.Role)
		}
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	// _, err = tx.Exec
	// return tx.Commit()
	return nil
}

func (r *RoleRepo) Update(ctx context.Context, role *proto_user.UpdateRoleRequest) error {
	query := fmt.Sprintf("UPDATE %s SET service=$1, role=$2 WHERE id=$3", r.tableName)

	_, err := r.db.Exec(query, role.Service, role.Role, role.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *RoleRepo) Delete(ctx context.Context, role *proto_user.DeleteRoleRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.tableName)

	_, err := r.db.Exec(query, role.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
