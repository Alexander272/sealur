package repo

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/config"
	"github.com/Alexander272/sealur/user_service/internal/models"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
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

func (r *IpRepo) GetAll(ctx context.Context, req *proto_user.GetAllIpRequest) (ips []models.Ip, err error) {
	query := fmt.Sprintf("SELECT ip, date, user_id FROM %s ORDER BY user_id", r.tableName)

	if err := r.db.Select(&ips, query); err != nil {
		return ips, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return ips, nil
}

func (r *IpRepo) Add(ctx context.Context, ip *proto_user.AddIpRequest) error {
	var count models.Count
	query := fmt.Sprintf("SELECT COUNT(ip) as count FROM %s GROUP BY user_id", r.tableName)
	if err := r.db.Get(&count, query); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	if count.Count >= r.conf.Max {
		query := fmt.Sprintf("DELETE FROM %s WHERE date = (SELECT min(date) FROM %s GROUP BY user_id)", r.tableName, r.tableName)
		_, err := r.db.Exec(query)
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
