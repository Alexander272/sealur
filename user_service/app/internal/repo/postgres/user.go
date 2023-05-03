package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
	query := fmt.Sprintf(`SELECT "%s".id, confirmed, company, inn, kpp, region, city, "position", phone, password, email, %s.code as role_code, name, address
		FROM "%s" INNER JOIN %s on %s.id=role_id WHERE email=$1`,
		UserTable, RoleTable, UserTable, RoleTable, RoleTable,
	)

	if err := r.db.Get(&data, query, req.Email); err != nil {
		return nil, "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	if !data.Confirmed {
		return nil, "", models.ErrUserNotVerified
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
		FROM "%s" INNER JOIN %s on %s.id=role_id WHERE %s.code='manager'`,
		UserTable, RoleTable, UserTable, RoleTable, RoleTable, RoleTable,
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

func (r *UserRepo) GetAnalytics(ctx context.Context) error {

	return fmt.Errorf("not implemented")
}

func (r *UserRepo) Create(ctx context.Context, user *user_api.CreateUser, roleId string) (string, error) {
	query := fmt.Sprintf(`INSERT INTO "%s"(id, company, inn, kpp, region, city, "position", phone, password, email, role_id, name, address, manager_id, use_link)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`, UserTable)
	id := uuid.New()

	managerId := user.ManagerId
	if managerId == "" {
		managerId = uuid.Nil.String()
	}

	_, err := r.db.Exec(query, id, user.Company, user.Inn, user.Kpp, user.Region, user.City, user.Position, user.Phone, user.Password,
		user.Email, roleId, user.Name, user.Address, managerId, user.UseLink,
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

func (r *UserRepo) Update(ctx context.Context, user *user_api.UpdateUser) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user.Company != "" {
		setValues = append(setValues, fmt.Sprintf("company=$%d", argId))
		args = append(args, user.Company)
		argId++
	}
	if user.Address != "" {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argId))
		args = append(args, user.Address)
		argId++
	}
	if user.Inn != "" {
		setValues = append(setValues, fmt.Sprintf("inn=$%d", argId))
		args = append(args, user.Inn)
		argId++
	}
	if user.Kpp != "" {
		setValues = append(setValues, fmt.Sprintf("kpp=$%d", argId))
		args = append(args, user.Kpp)
		argId++
	}
	if user.Region != "" {
		setValues = append(setValues, fmt.Sprintf("region=$%d", argId))
		args = append(args, user.Region)
		argId++
	}
	if user.City != "" {
		setValues = append(setValues, fmt.Sprintf("city=$%d", argId))
		args = append(args, user.City)
		argId++
	}
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
	if user.Password != "" {
		setValues = append(setValues, fmt.Sprintf("password=$%d", argId))
		args = append(args, user.Password)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE "%s" SET %s WHERE id=$%d`, UserTable, setQuery, argId)

	args = append(args, user.Id)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
