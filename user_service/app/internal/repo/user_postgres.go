package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/user_service/internal/models"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
	"github.com/google/uuid"
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
		query = fmt.Sprintf("SELECT id, email, password FROM %s WHERE login = $1 AND confirmed=true", r.tableName)
		param = req.Login
	} else {
		query = fmt.Sprintf("SELECT id, organization, name, email, city, position, phone FROM %s WHERE id = $1", r.tableName)
		param = req.Id
	}

	if err := r.db.Get(&user, query, param); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, err
		}
		return user, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return user, nil
}

func (r *UserRepo) GetAll(ctx context.Context, req *proto_user.GetAllUserRequest) (users []models.User, err error) {
	query := fmt.Sprintf(`SELECT id, organization, name, email, city, position, phone, login FROM %s WHERE confirmed=true 
		ORDER BY organization, name, id`, r.tableName)

	if err := r.db.Select(&users, query); err != nil {
		return users, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return users, nil
}

func (r *UserRepo) GetNew(ctx context.Context, req *proto_user.GetNewUserRequest) (users []models.User, err error) {
	query := fmt.Sprintf("SELECT id, organization, name, email, city, position, phone FROM %s WHERE confirmed=false ORDER BY date_reg", r.tableName)

	if err := r.db.Select(&users, query); err != nil {
		return users, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, user *proto_user.CreateUserRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, organization, name, email, city, "position", phone) VALUES ($1, $2, $3, $4, $5, $6, $7)`, r.tableName)
	id := uuid.New()

	_, err := r.db.Exec(query, id, user.Organization, user.Name, user.Email, user.City, user.Position, user.Phone)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *UserRepo) Confirm(ctx context.Context, user *proto_user.ConfirmUserRequest) (models.ConfirmUser, error) {
	query := fmt.Sprintf("UPDATE %s SET login=$1, password=$2, confirmed=true WHERE id=$3 RETURNING name, email", r.tableName)
	row := r.db.QueryRow(query, user.Login, user.Password, user.Id)

	if row.Err() != nil {
		return models.ConfirmUser{}, fmt.Errorf("failed to execute query. error: %w", row.Err())
	}

	var u models.ConfirmUser
	if err := row.Scan(&u.Name, &u.Email); err != nil {
		return models.ConfirmUser{}, fmt.Errorf("failed to scan result. error: %w", err)
	}

	return u, nil
}

func (r *UserRepo) Update(ctx context.Context, user *proto_user.UpdateUserRequest) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, user.Name)
		argId++
	}
	if user.Email != "" {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, user.Email)
		argId++
	}
	if user.Position != "" {
		setValues = append(setValues, fmt.Sprintf("position=$%d", argId))
		args = append(args, user.Position)
		argId++
	}
	if user.Phone != "" {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argId))
		args = append(args, user.Phone)
		argId++
	}
	if user.Login != "" {
		setValues = append(setValues, fmt.Sprintf("login=$%d", argId))
		args = append(args, user.Login)
		argId++
	}
	if user.Password != "" {
		setValues = append(setValues, fmt.Sprintf("password=$%d", argId))
		args = append(args, user.Password)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", r.tableName, setQuery, argId)

	args = append(args, user.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, user *proto_user.DeleteUserRequest) (models.DeleteUser, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING name, email", r.tableName)
	row := r.db.QueryRow(query, user.Id)

	if row.Err() != nil {
		return models.DeleteUser{}, fmt.Errorf("failed to execute query. error: %w", row.Err())
	}

	var u models.DeleteUser
	if err := row.Scan(&u.Name, &u.Email); err != nil {
		return models.DeleteUser{}, fmt.Errorf("failed to scan result. error: %w", err)
	}

	return u, nil
	// _, err := r.db.Exec(query, user.Id)
	// if err != nil {
	// 	return fmt.Errorf("failed to execute query. error: %w", err)
	// }

	// return nil
}
