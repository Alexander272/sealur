package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

// Получение данных
func (s *DataService) GetData(ctx context.Context, data *calc_api.FlangeRequest) (result models.DataFlange, err error) {
	// Получение данных о прокладке
	result.Gasket, result.TypeGasket, err = s.getGasketData(ctx, data.Gasket)
	if err != nil {
		return result, err
	}

	// Получение данных о фланце и болтах
	flange1, boltSize, err := s.getFlangeData(ctx, data.FlangesData[0], data.Bolts, data.Flanges.String(), data.Temp)
	if err != nil {
		return result, err
	}

	result.Type1, result.Type2 = data.FlangesData[0].Type, data.FlangesData[0].Type
	flange2 := flange1

	if len(data.FlangesData) > 1 {
		flange2, _, err = s.getFlangeData(ctx, data.FlangesData[1], data.Bolts, data.Flanges.String(), data.Temp)
		if err != nil {
			return result, err
		}
		result.Type2 = data.FlangesData[1].Type
	}
	result.Flange1 = flange1
	result.Flange2 = flange2

	if data.IsEmbedded {
		// Получение данных для закладной детали
		result.Embed, err = s.getEmbedData(ctx, data.Embed, data.Temp)
		if err != nil {
			return result, err
		}
	}

	//* формула из Таблицы В.1
	Tb := s.typeFlangesTB[data.Flanges.String()] * data.Temp
	if data.FlangesData[0].Type == flange_model.FlangeData_free {
		Tb = s.typeFlangesTB[data.Flanges.String()+"-free"] * data.Temp
	}

	Lb0 := result.Gasket.Thickness + result.Flange1.H + result.Flange2.H
	if result.Type1 == flange_model.FlangeData_free {
		Lb0 += result.Flange1.Hk
	}
	if result.Type2 == flange_model.FlangeData_free {
		Lb0 += result.Flange2.Hk
	}
	if data.IsEmbedded {
		Lb0 += result.Gasket.Thickness + result.Embed.Thickness
	}

	// получение данных о болте (инфа о материале, также размеры записываются в объект)
	result.Bolt, err = s.getBoltData(ctx, data.Bolts, boltSize, Lb0, Tb)
	if err != nil {
		return result, err
	}

	if data.IsUseWasher {
		// получение данных о шайбах
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

	return result, nil
}
