package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (r *DeviceRepo) GetNameGasket(ctx context.Context, req *device_api.GetNameGasketRequest) (gasket []models.NameGasketDTO, err error) {
	query := fmt.Sprintf("SELECT id, num_id, pres_id, title FROM %s WHERE fin_id=$1", NumberOfMovesTable)

	if err := r.db.Get(&gasket, query, req.FinId); err != nil {
		return gasket, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return gasket, nil
}

func (r *DeviceRepo) CreateNameGasket(ctx context.Context, gasket *device_api.CreateNameGasketRequest) (id string, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (fin_id, num_id, pres_id, title, size_long, size_trans, width, thick1, thick2, thick3, thick4) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`, NameGasketTable)

	row := r.db.QueryRow(query, gasket.FinId, gasket.NumId, gasket.PresId, gasket.Title, gasket.SizeLong, gasket.SizeTrans, gasket.Width,
		gasket.Thick1, gasket.Thick2, gasket.Thick3, gasket.Thick4)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	var returnedId int
	if err := row.Scan(&returnedId); err != nil {
		return "", fmt.Errorf("failed to scan result. error: %w", err)
	}

	return fmt.Sprintf("%d", returnedId), nil
}

func (r *DeviceRepo) CreateFewNameGasket(ctx context.Context, gasket *device_api.CreateFewNameGasketRequest) error {
	query := fmt.Sprintf("INSERT INTO %s (fin_id, num_id, pres_id, title, size_long, size_trans, width, thick1, thick2, thick3, thick4) VALUES ",
		NameGasketTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(gasket.Gasket))

	c := 11
	for i, d := range gasket.Gasket {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*c+1, i*c+2, i*c+3, i*c+4, i*c+5, i*c+6, i*c+7, i*c+8, i*c+9, i*c+10, i*c+11))
		args = append(args, d.FinId, d.NumId, d.PresId, d.Title, d.SizeLong, d.SizeTrans, d.Width, d.Thick1, d.Thick2, d.Thick3, d.Thick4)
	}
	query += strings.Join(values, ", ")

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *DeviceRepo) UpdateNameGasket(ctx context.Context, gasket *device_api.UpdateNameGasketRequest) error {
	query := fmt.Sprintf(`UPDATE %s SET fin_id=$1, num_id=$2, pres_id=$3, title=$4, size_long=$5, size_trans=$6, width=$7, 
		thick1=$8, thick2=$9, thick3=$10, thick4=$11 WHERE id=$12`, NameGasketTable)

	_, err := r.db.Exec(query, gasket.FinId, gasket.NumId, gasket.PresId, gasket.Title, gasket.SizeLong, gasket.SizeTrans, gasket.Width,
		gasket.Thick1, gasket.Thick2, gasket.Thick3, gasket.Thick4)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *DeviceRepo) DeleteNameGasket(ctx context.Context, gasket *device_api.DeleteNameGasketRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", NameGasketTable)

	if _, err := r.db.Exec(query, gasket.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
