package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
	"github.com/jmoiron/sqlx"
)

type DeviceRepo struct {
	db *sqlx.DB
}

func NewDeviceRepo(db *sqlx.DB) *DeviceRepo {
	return &DeviceRepo{db: db}
}

func (r *DeviceRepo) GetDevices(ctx context.Context, req *device_api.GetDeviceRequest) (devices []models.DeviceDTO, err error) {
	query := fmt.Sprintf("SELECT id, title FROM %s", DeviceModTable)

	if err := r.db.Get(&devices, query); err != nil {
		return devices, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return devices, nil
}

func (r *DeviceRepo) CreateDevices(ctx context.Context, device *device_api.CreateDeviceRequest) (id string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", DeviceModTable)

	row := r.db.QueryRow(query, device.Title)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var returnedId int
	if err := row.Scan(&returnedId); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", returnedId), nil
}

func (r *DeviceRepo) CreateFewDevices(ctx context.Context, device *device_api.CreateFewDeviceRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ", DeviceModTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(device.Divices))

	for i, d := range device.Divices {
		values = append(values, fmt.Sprintf("($%d)", i+1))
		args = append(args, d.Title)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *DeviceRepo) UpdateDevice(ctx context.Context, device *device_api.UpdateDeviceRequest) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1 WHERE id=$2", DeviceModTable)

	_, err := r.db.Exec(query, device.Title, device.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *DeviceRepo) DeleteDevice(ctx context.Context, device *device_api.DeleteDeviceRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", DeviceModTable)

	if _, err := r.db.Exec(query, device.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
