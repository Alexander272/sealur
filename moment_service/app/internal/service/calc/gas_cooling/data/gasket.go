package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/gas_cooling_model"
	"github.com/Alexander272/sealur_proto/api/moment/device_api"
)

func (s *DataService) getGasketData(ctx context.Context, gasket *gas_cooling_model.GasketData,
) (gasketData *gas_cooling_model.GasketResult, typeGasket gas_cooling_model.GasketData_Type, err error) {

	dto := models.GetGasket{GasketId: gasket.GasketId, EnvId: gasket.EnvId, Thickness: gasket.Thickness}
	g, err := s.gasket.GetFullData(ctx, dto)
	if err != nil {
		return nil, 0, err
	}

	nameGasket, err := s.device.GetNameGasketSize(ctx, &device_api.GetNameGasketSizeRequest{Id: gasket.NameGasket.Id})
	if err != nil {
		return nil, 0, err
	}

	res := &gas_cooling_model.GasketResult{
		Gasket:          g.Gasket,
		Env:             g.Env,
		Type:            g.Type,
		Thickness:       gasket.Thickness,
		SizeLong:        nameGasket.SizeLong,
		SizeTrans:       nameGasket.SizeTrans,
		Width:           nameGasket.Width,
		Thick1:          nameGasket.Thick1,
		Thick2:          nameGasket.Thick2,
		Thick3:          nameGasket.Thick3,
		Thick4:          nameGasket.Thick4,
		M:               g.M,
		Pres:            g.SpecificPres,
		PermissiblePres: g.PermissiblePres,
		Compression:     g.Compression,
		Epsilon:         g.Epsilon,
		Name:            gasket.NameGasket.Title,
	}
	return res, gas_cooling_model.GasketData_Type(gas_cooling_model.GasketData_Type_value[g.Type]), nil

}
