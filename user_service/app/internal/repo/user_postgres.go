package repo

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/models"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db        *sqlx.DB
	tableName string
}

func NewUserRepo(db *sqlx.DB, tableName string) *UserRepo {
	return &UserRepo{db: db, tableName: tableName}
}

func (r *UserRepo) Get(ctx context.Context, req *proto_user.GetUserRequest) (user models.User, err error) {
	var query, param string
	if req.Login != "" {
		query = fmt.Sprintf("SELECT id, email, password FROM %s WHERE login = $1", r.tableName)
		param = req.Login
	} else {
		query = fmt.Sprintf("SELECT id, organization, name, email, city, position, phone FROM %s WHERE id = $1", r.tableName)
		param = req.Id
	}

	if err := r.db.Get(&user, query, param); err != nil {
		return user, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return user, nil
}

func (r *UserRepo) GetAll(ctx context.Context, req *proto_user.GetAllUserRequest) (users []models.User, err error) {
	query := fmt.Sprintf("SELECT id, organization, name, email, city, position, phone FROM %s WHERE confirmed=true ORDER BY id", r.tableName)

	if err := r.db.Select(&users, query); err != nil {
		return users, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return users, nil
}

func (r *UserRepo) GetNew(ctx context.Context, req *proto_user.GetNewUserRequest) (users []models.User, err error) {
	query := fmt.Sprintf("SELECT id, organization, name, email, city, position, phone FROM %s WHERE confirmed=false", r.tableName)

	if err := r.db.Select(&users, query); err != nil {
		return users, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, user *proto_user.CreateUserRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (organization, name, email, city, "position", phone) VALUES ($1, $2, $3, $4, $5, $6)`, r.tableName)

	_, err := r.db.Exec(query, user.Organization, user.Name, user.Email, user.City, user.Position, user.Phone)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *UserRepo) Confirm(ctx context.Context, user *proto_user.ConfirmUserRequest) error {
	query := fmt.Sprintf("UPDATE %s SET login=$1, password=$2 WHERE id=$3", r.tableName)

	_, err := r.db.Exec(query, user.Login, user.Password, user.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *UserRepo) Update(ctx context.Context, user *proto_user.UpdateUserRequest) error {
	query := fmt.Sprintf("UPDATE %s SET name=$1, email=$2, position=$3, phone=$4 WHERE id=$5", r.tableName)

	_, err := r.db.Exec(query, user.Name, user.Email, user.Position, user.Phone, user.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, user *proto_user.DeleteUserRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.tableName)

	_, err := r.db.Exec(query, user.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
