package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/config"
	"github.com/Alexander272/sealur/user_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/user_api"
	"github.com/jmoiron/sqlx"
)

type IpRepo struct {
	db        *sqlx.DB
	tableName string
	conf      config.IPConfig
}

func NewIpRepo(db *sqlx.DB, tableName string, conf config.IPConfig) *IpRepo {
	return &IpRepo{
		db:        db,
		tableName: tableName,
		conf:      conf,
	}
}

func (r *IpRepo) GetAll(ctx context.Context, req *user_api.GetAllIpRequest) (ips []models.Ip, err error) {
	query := fmt.Sprintf("SELECT ip, date, user_id FROM %s ORDER BY user_id", r.tableName)

	if err := r.db.Select(&ips, query); err != nil {
		return ips, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return ips, nil
}

func (r *IpRepo) Add(ctx context.Context, ip *user_api.AddIpRequest) error {
	var count models.Count
	query := fmt.Sprintf("SELECT COUNT(ip) as count FROM %s GROUP BY user_id HAVING user_id=$1", r.tableName)
	if err := r.db.Get(&count, query, ip.UserId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			count.Count = 0
		} else {
			return fmt.Errorf("failed to execute query. error: %w", err)
		}
	}

	if count.Count >= r.conf.Max {
		query := fmt.Sprintf("DELETE FROM %s WHERE date = (SELECT min(date) FROM %s GROUP BY user_id HAVING user_id=$1)", r.tableName, r.tableName)
		_, err := r.db.Exec(query, ip.UserId)
		if err != nil {
			return fmt.Errorf("failed to execute query. error: %w", err)
		}
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, ip) VALUES ($1, $2)", r.tableName)
	_, err := r.db.Exec(query, ip.UserId, ip.Ip)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
