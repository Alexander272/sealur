package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Alexander272/sealur/user_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/user/models/user_model"
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

func (r *UserRepo) Get(ctx context.Context, req *user_api.GetUser) (*user_model.User, error) {
	var data models.User
	query := fmt.Sprintf(`SELECT "%s".id, company, inn, kpp, region, city, "position", phone, password, email, %s.code as role_code, name, address, manager_id
		FROM "%s" INNER JOIN %s on %s.id=role_id WHERE "%s".id=$1`,
		UserTable, RoleTable, UserTable, RoleTable, RoleTable, UserTable,
	)

	if err := r.db.Get(&data, query, req.Id); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	user := &user_model.User{
		Id:        data.Id,
		Company:   data.Company,
		Inn:       data.Inn,
		Kpp:       data.Kpp,
		Region:    data.Region,
		City:      data.City,
		Position:  data.Position,
		Phone:     data.Phone,
		Email:     data.Email,
		RoleCode:  data.RoleCode,
		Name:      data.Name,
		Address:   data.Address,
		ManagerId: data.ManagerId,
	}

	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, req *user_api.GetUserByEmail) (*user_model.User, string, error) {
	var data models.User
	query := fmt.Sprintf(`SELECT "%s".id, company, inn, kpp, region, city, "position", phone, password, email, %s.code as role_code, name, address
		FROM "%s" INNER JOIN %s on %s.id=role_id WHERE email=$1 AND confirmed=true`,
		UserTable, RoleTable, UserTable, RoleTable, RoleTable,
	)

	if err := r.db.Get(&data, query, req.Email); err != nil {
		return nil, "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	user := &user_model.User{
		Id:       data.Id,
		Company:  data.Company,
		Inn:      data.Inn,
		Kpp:      data.Kpp,
		Region:   data.Region,
		City:     data.City,
		Position: data.Position,
		Phone:    data.Phone,
		Email:    data.Email,
		RoleCode: data.RoleCode,
		Name:     data.Name,
		Address:  data.Address,
	}

	return user, data.Password, nil
}

func (r *UserRepo) GetManagers(ctx context.Context, req *user_api.GetNewUser) (users []*user_model.User, err error) {
	var data []models.User
	query := fmt.Sprintf(`SELECT "%s".id, region, city, "position", phone, email, %s.code as role_code, name
		FROM "%s" INNER JOIN %s on %s.id=role_id WHERE role_code=manager`,
		UserTable, RoleTable, UserTable, RoleTable, RoleTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, u := range data {
		users = append(users, &user_model.User{
			Id:       u.Id,
			Region:   u.Region,
			City:     u.City,
			Position: u.Position,
			Phone:    u.Phone,
			Email:    u.Email,
			Name:     u.Name,
			RoleCode: u.RoleCode,
		})
	}
	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, user *user_api.CreateUser, roleId string) (string, error) {
	query := fmt.Sprintf(`INSERT INTO "%s"(id, company, inn, kpp, region, city, "position", phone, password, email, role_id, name, address, manager_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, UserTable)
	id := uuid.New()

	managerId := user.ManagerId
	if managerId == "" {
		managerId = uuid.Nil.String()
	}

	_, err := r.db.Exec(query, id, user.Company, user.Inn, user.Kpp, user.Region, user.City, user.Position, user.Phone, user.Password,
		user.Email, roleId, user.Name, user.Address, managerId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id.String(), nil
}

func (r *UserRepo) Confirm(ctx context.Context, user *user_api.ConfirmUser) error {
	query := fmt.Sprintf(`UPDATE "%s" SET confirmed=true, date=$1 WHERE id=$2`, UserTable)
	date := fmt.Sprintf("%d", time.Now().UnixMilli())

	_, err := r.db.Exec(query, date, user.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *UserRepo) SetManager(ctx context.Context, manager *user_api.UserManager) error {
	query := fmt.Sprintf(`UPDATE "%s" SET manager_id=$1 WHERE id=$2`, UserTable)

	_, err := r.db.Exec(query, manager.ManagerId, manager.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
