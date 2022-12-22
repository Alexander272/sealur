package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
)

func (s *DataService) GetData(ctx context.Context, data *calc_api.DevCoolingRequest) (result models.DataDevCooling, err error) {
	result.Cap, err = s.getDataCap(ctx, data.Cap, data.Temp)
	if err != nil {
		return result, err
	}

	result.TubeSheet, err = s.getTubeSheetData(ctx, data.TubeSheet, data.Temp)
	if err != nil {
		return result, err
	}

	result.Tube, err = s.getTubeData(ctx, data.Tube, data.Temp)
	if err != nil {
		return result, err
	}

	result.Bolt, err = s.getBoltData(ctx, data.Bolts, data.Temp)
	if err != nil {
		return result, err
	}

	result.Gasket, result.TypeGasket, err = s.getGasketData(ctx, data.Gasket)
	if err != nil {
		return result, err
	}

	return result, nil
}
