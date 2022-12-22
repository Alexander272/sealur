package data

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

func (s *DataService) GetDataOld(ctx context.Context, data *calc_api.FlangeRequest) (result models.DataFlangeOld, err error) {
	//* формула из Таблицы В.1
	Tb := s.typeFlangesTB[data.Flanges.String()] * data.Temp
	if data.FlangesData[0].Type == flange_model.FlangeData_free {
		Tb = s.typeFlangesTB[data.Flanges.String()+"-free"] * data.Temp
	}

	flange1, boltSize, err := s.getDataFlange(ctx, data.FlangesData[0], data.Bolts, data.Flanges.String(), data.Temp)
	if err != nil {
		return result, err
	}

	result.Type1, result.Type2 = data.FlangesData[0].Type, data.FlangesData[0].Type
	flange2 := flange1

	if len(data.FlangesData) > 1 {
		flange2, _, err = s.getDataFlange(ctx, data.FlangesData[1], data.Bolts, data.Flanges.String(), data.Temp)
		if err != nil {
			return result, err
		}
		result.Type2 = data.FlangesData[1].Type
	}

	result.Bolt, err = s.getBoltData(ctx, data.Bolts, boltSize, flange1.L, Tb)
	if err != nil {
		return result, err
	}

	if data.IsUseWasher {
		result.Washer1, err = s.getWasherData(ctx, data.Washer[0], flange1.Tf)
		if err != nil {
			return result, err
		}
		if !data.IsSameFlange {
			result.Washer2, err = s.getWasherData(ctx, data.Washer[1], flange2.Tf)
			if err != nil {
				return result, err
			}
		} else {
			result.Washer2 = result.Washer1
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

	flange1, err = s.getCalculatedDataFlange(ctx, data.FlangesData[0].Type, flange1, result.Dcp)
	if err != nil {
		return result, err
	}

	if len(data.FlangesData) > 1 {
		flange2, err = s.getCalculatedDataFlange(ctx, data.FlangesData[1].Type, flange2, result.Dcp)
		if err != nil {
			return result, err
		}
	} else {
		flange2 = flange1
	}

	result.Flange1 = flange1
	result.Flange2 = flange2

	return result, nil
}
