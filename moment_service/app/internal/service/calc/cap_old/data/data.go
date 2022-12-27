package data

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

func (s *DataService) GetData(ctx context.Context, data *calc_api.CapRequestOld) (result models.DataCapOld, err error) {
	//* формула из Таблицы В.1
	Tb := s.typeFlangesTB[data.Flanges.String()] * data.Temp
	if data.FlangeData.Type == cap_model.FlangeData_free {
		Tb = s.typeFlangesTB[data.Flanges.String()+"-free"] * data.Temp
	}

	flange, boltSize, err := s.getDataFlange(ctx, data.FlangeData, data.Bolts, data.Flanges.String(), data.Temp)
	if err != nil {
		return result, err
	}
	cap, err := s.getDataCap(ctx, data.CapData, data.Flanges.String(), data.Temp)
	if err != nil {
		return result, err
	}

	result.FType = data.FlangeData.Type
	result.CType = data.CapData.Type

	result.Bolt, err = s.getBoltData(ctx, data.Bolts, boltSize, flange.L, Tb)
	if err != nil {
		return result, err
	}

	//? я использую температуру фланца. хз верно или нет.
	if data.IsUseWasher {
		result.Washer1, err = s.getWasherData(ctx, data.Washer[0], flange.Tf)
		if err != nil {
			return result, err
		}

		result.Washer2, err = s.getWasherData(ctx, data.Washer[1], cap.T)
		if err != nil {
			return result, err
		}
	}
	if data.IsEmbedded {
		result.Embed, err = s.getEmbedData(ctx, data.Embed, data.Temp)
		if err != nil {
			return result, err
		}
	}
	bp := (data.Gasket.DOut - data.Gasket.DIn) / 2
	result.Gasket, result.TypeGasket, err = s.getGasketData(ctx, data.Gasket, bp)
	if err != nil {
		return result, err
	}

	if result.TypeGasket == "Oval" {
		// фомула 4
		result.B0 = bp / 4
		// фомула ?
		result.Dcp = data.Gasket.DOut - bp/2

	} else {
		if bp <= constants.Bp {
			// фомула 2
			result.B0 = bp
		} else {
			// фомула 3
			result.B0 = constants.B0 * math.Sqrt(bp)
		}
		// фомула 5
		result.Dcp = data.Gasket.DOut - result.B0
	}

	flange, err = s.getCalculatedDataFlange(ctx, data.FlangeData.Type, flange, result.Dcp)
	if err != nil {
		return result, err
	}
	cap, err = s.getCalculatedDataCap(ctx, data.CapData.Type, cap, flange.H, flange.D, flange.S0, flange.DOut, result.Dcp)
	if err != nil {
		return result, err
	}

	result.Flange = flange
	result.Cap = cap

	return result, nil
}
