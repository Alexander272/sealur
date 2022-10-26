package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
)

func (s *DataService) GetData(ctx context.Context, data *calc_api.ExpressCircleRequest) (result models.DataExCircle, err error) {
	result.Bolt, err = s.getBoltData(ctx, data.Bolts, 20)
	if err != nil {
		return result, err
	}

	result.Gasket, result.TypeGasket, err = s.getGasketData(ctx, data.Gasket)
	if err != nil {
		return result, err
	}

	return result, nil
}
