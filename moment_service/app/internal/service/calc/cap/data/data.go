package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

func (s *DataService) GetData(ctx context.Context, data *calc_api.CapRequest) (result models.DataCap, err error) {
	//* формула из Таблицы В.1
	Tb := s.typeFlangesTB[data.Data.Flanges.String()] * data.Data.Temp
	if data.FlangeData.Type == cap_model.FlangeData_free {
		Tb = s.typeFlangesTB[data.Data.Flanges.String()+"-free"] * data.Data.Temp
	}

	flange, boltSize, err := s.getDataFlange(ctx, data.FlangeData, data.Bolts, data.Data.Flanges.String(), data.Data.Temp)
	if err != nil {
		return result, err
	}
	cap, err := s.getDataCap(ctx, data.CapData, data.Data.Flanges.String(), data.Data.Temp)
	if err != nil {
		return result, err
	}

	result.FlangeType = data.FlangeData.Type
	result.CapType = data.CapData.Type

	result.Bolt, err = s.getBoltData(ctx, data.Bolts, boltSize, flange.L, Tb)
	if err != nil {
		return result, err
	}

	//? я использую температуру фланца. хз верно или нет.
	if data.IsUseWasher {
		result.Washer1, err = s.getWasherData(ctx, data.Washer[0], flange.T)
		if err != nil {
			return result, err
		}

		result.Washer2, err = s.getWasherData(ctx, data.Washer[1], cap.T)
		if err != nil {
			return result, err
		}
	}
	if data.Data.IsEmbedded {
		result.Embed, err = s.getEmbedData(ctx, data.Embed, data.Data.Temp)
		if err != nil {
			return result, err
		}
	}
	bp := (data.Gasket.DOut - data.Gasket.DIn) / 2
	result.Gasket, result.TypeGasket, err = s.getGasketData(ctx, data.Gasket, bp)
	if err != nil {
		return result, err
	}

	result.Flange = flange
	result.Cap = cap

	return result, nil
}
