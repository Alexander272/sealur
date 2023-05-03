package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_configuration_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_conf_api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PutgConfigurationRepo struct {
	db *sqlx.DB
}

func NewPutgConfigurationRepo(db *sqlx.DB) *PutgConfigurationRepo {
	return &PutgConfigurationRepo{
		db: db,
	}
}

func (r *PutgConfigurationRepo) Get(ctx context.Context, req *putg_conf_api.GetPutgConfiguration,
) (configuration []*putg_configuration_model.PutgConfiguration, err error) {
	var data []models.PutgConfiguration
	query := fmt.Sprintf(`SELECT id, title, code, has_standard, has_drawing FROM %s ORDER BY is_default DESC`, PutgConfTable)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pc := range data {
		configuration = append(configuration, &putg_configuration_model.PutgConfiguration{
			Id:          pc.Id,
			Title:       pc.Title,
			Code:        pc.Code,
			HasStandard: pc.HasStandard,
			HasDrawing:  pc.HasDrawing,
		})
	}

	return configuration, nil
}

func (r *PutgConfigurationRepo) Create(ctx context.Context, c *putg_conf_api.CreatePutgConfiguration) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, title, code, has_standard, has_drawing, is_default)
		VALUES ($1, $2, $3, $4, $5, $6)`, PutgConfTable,
	)
	id := uuid.New()

	_, err := r.db.Exec(query, id, c.Title, c.Code, c.HasStandard, c.HasDrawing, c.IsDefault)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgConfigurationRepo) Update(ctx context.Context, c *putg_conf_api.UpdatePutgConfiguration) error {
	query := fmt.Sprintf(`UPDATE %s SET title=$1, code=$2, has_standard=$3, has_drawing=$4, is_default=$5 WHERE id=$6`, PutgConfTable)

	_, err := r.db.Exec(query, c.Title, c.Code, c.HasStandard, c.HasDrawing, c.IsDefault, c.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *PutgConfigurationRepo) Delete(ctx context.Context, c *putg_conf_api.DeletePutgConfiguration) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PutgConfTable)

	_, err := r.db.Exec(query, c.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}
