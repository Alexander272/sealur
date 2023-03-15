package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// TODO
func (r *UserRepo) GetByEmail(ctx context.Context) {}

// TODO
func (r *UserRepo) Get(ctx context.Context) {}

func (r *UserRepo) Create(ctx context.Context, user *user_api.CreateUser, roleId string) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, company, inn, kpp, region, city, "position", phone, password, email, role_id, name, address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`, UserTable)
	id := uuid.New()

	_, err := r.db.Exec(query, id, user.Company, user.Inn, user.Kpp, user.Region, user.City, user.Position, user.Phone, user.Password,
		user.Email, roleId, user.Name, user.Address,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

// TODO возможно нужно возвращать пользователя при подтверждении
func (r *UserRepo) Confirm(ctx context.Context, user *user_api.ConfirmUser) error {
	query := fmt.Sprintf(`UPDATE %s	SET confirmed=true, date=$1 WHERE id=$2`, UserTable)
	date := fmt.Sprintf("%d", time.Now().UnixMilli())

	_, err := r.db.Exec(query, date, user.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
