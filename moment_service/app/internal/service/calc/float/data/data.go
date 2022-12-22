package data

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

func (s *DataService) GetData(ctx context.Context, data *calc_api.FloatRequest) (result models.DataFloat, err error) {
	flange, boltSize, err := s.getDataFlange(ctx, data.FlangeData, data.Bolts)
	if err != nil {
		return result, err
	}
	cap, err := s.getDataCap(ctx, data.CapData)
	if err != nil {
		return result, err
	}

	result.Bolt, err = s.getBoltData(ctx, data.Bolts, boltSize, data.Bolts.Distance, data.Bolts.Temp)
	if err != nil {
		return result, err
	}

	bp := (data.Gasket.DOut - data.Gasket.DIn) / 2

	result.Gasket, result.TypeGasket, err = s.getGasketData(ctx, data.Gasket, bp)
	if err != nil {
		return result, err
	}

	if data.HasThorn {
		result.B0 = (flange.Width + bp) / 2
		result.Dcp = flange.DIn + flange.Width
	} else {
		if result.TypeGasket != float_model.GasketData_Soft.String() {
			result.B0 = bp / 4
			result.Dcp = data.Gasket.DOut - bp/2
		} else {
			if bp <= constants.Bp {
				result.B0 = bp
			} else {
				result.B0 = constants.B0 * math.Sqrt(bp)
			}
			result.Dcp = data.Gasket.DOut - result.B0
		}
	}

	flange, err = s.getCalculatedDataFlange(ctx, flange, cap, result.Dcp)
	if err != nil {
		return result, err
	}
	cap, err = s.getCalculatedDataCap(ctx, flange, cap)
	if err != nil {
		return result, err
	}

	result.Flange = flange
	result.Cap = cap

	return result, nil
}
