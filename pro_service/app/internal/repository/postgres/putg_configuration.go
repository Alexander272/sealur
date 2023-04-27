package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_configuration_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_conf_api"
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

// TODO дописать оставшиеся функции
