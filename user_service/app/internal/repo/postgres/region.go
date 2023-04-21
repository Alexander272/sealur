package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/models"
	"github.com/Alexander272/sealur/user_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/user/models/user_model"
	"github.com/jmoiron/sqlx"
)

type RegionRepo struct {
	db *sqlx.DB
}

func NewRegionRepo(db *sqlx.DB) *RegionRepo {
	return &RegionRepo{
		db: db,
	}
}

// func (r *RegionRepo) GetAll(ctx context.Context)

func (r *RegionRepo) GetManagerByRegion(ctx context.Context, region string) (*user_model.User, error) {
	var data models.User
	query := fmt.Sprintf(`SELECT %s.manager_id as id, company, inn, kpp, region, city, "position", phone, email, name, address
		FROM %s INNER JOIN "%s" ON "%s".id=%s.manager_id WHERE title=$1`, RegionTable, RegionTable, UserTable, UserTable, RegionTable,
	)

	if err := r.db.Get(&data, query, region); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	query = fmt.Sprintf(`UPDATE %s SET count_query=count_query+1 WHERE title=$1`, RegionTable)
	_, err := r.db.Exec(query, region)
	if err != nil {
		logger.Error("failed to update count_query. error: %w", err)
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
		Name:     data.Name,
		Address:  data.Address,
	}

	return user, nil
}
