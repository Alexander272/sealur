package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/putg_data_model"
	"github.com/Alexander272/sealur_proto/api/pro/putg_data_api"
	"github.com/jmoiron/sqlx"
)

type PutgDataRepo struct {
	db *sqlx.DB
}

func NewPutgDataRepo(db *sqlx.DB) *PutgDataRepo {
	return &PutgDataRepo{
		db: db,
	}
}

// TODO эта функция косячная. надо определиться как я буду получать данные
// func (r *PutgDataRepo) Get(ctx context.Context, req *putg_data_api.GetPutgData) (putg []*putg_data_model.PutgData, err error) {
// 	var data []models.PutgData
// 	query := fmt.Sprintf(`SELECT id, filler_id, has_jumper, has_hole, has_removable, has_mounting, has_coating WHERE filler_id=$1 FROM %s`, PutgDataTable)

// 	if err := r.db.Select(&data, query, req.FillerId); err != nil {
// 		return nil, fmt.Errorf("failed to execute query. error: %w", err)
// 	}

// 	for _, pd := range data {
// 		putg = append(putg, &putg_data_model.PutgData{
// 			Id:           pd.Id,
// 			HasJumper:    pd.HasJumper,
// 			HasHole:      pd.HasHole,
// 			HasRemovable: pd.HasRemovable,
// 			HasMounting:  pd.HasMounting,
// 			HasCoating:   pd.HasCoating,
// 		})
// 	}

// 	return putg, nil
// }

// ? хз правильно ли это работает (надо тестировать)
func (r *PutgDataRepo) Get(ctx context.Context, req *putg_data_api.GetPutgData) (putg []*putg_data_model.PutgData, err error) {
	var data []models.PutgData
	// query := fmt.Sprintf(`SELECT id, filler_id, has_jumper, has_hole, has_removable, has_mounting, has_coating FROM %s WHERE filler_id in (
	// 	SELECT id FROM %s WHERE construction_id=$1 ORDER BY code)`, PutgDataTable, PutgFillerTable,
	// )
	query := fmt.Sprintf(`SELECT %s.id, filler_id, has_jumper, has_hole, has_removable, has_mounting, has_coating FROM %s 
		INNER JOIN %s ON filler_id=%s.id WHERE construction_id=$1 ORDER BY code`,
		PutgDataTable, PutgDataTable, PutgFillerTable, PutgFillerTable,
	)

	if err := r.db.Select(&data, query, req.ConstructionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, pd := range data {
		putg = append(putg, &putg_data_model.PutgData{
			Id:           pd.Id,
			FillerId:     pd.FillerId,
			HasJumper:    pd.HasJumper,
			HasHole:      pd.HasHole,
			HasRemovable: pd.HasRemovable,
			HasMounting:  pd.HasMounting,
			HasCoating:   pd.HasCoating,
		})
	}

	return putg, nil
}
