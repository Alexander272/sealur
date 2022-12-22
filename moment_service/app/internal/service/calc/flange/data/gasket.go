package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

func (s *DataService) getGasketData(ctx context.Context, data *flange_model.GasketData,
) (*flange_model.GasketResult, flange_model.GasketData_Type, error) {
	bp := (data.DOut - data.DIn) / 2

	if data.GasketId != "another" {
		g := models.GetGasket{GasketId: data.GasketId, EnvId: data.EnvId, Thickness: data.Thickness}
		gasket, err := s.gasket.GetFullData(ctx, g)
		if err != nil {
			return nil, 0, err
		}

		res := &flange_model.GasketResult{
			Gasket:          gasket.Gasket,
			Env:             gasket.Env,
			Type:            gasket.Type,
			Thickness:       data.Thickness,
			DOut:            data.DOut,
			Width:           bp,
			M:               gasket.M,
			Pres:            gasket.SpecificPres,
			PermissiblePres: gasket.PermissiblePres,
			Compression:     gasket.Compression,
			Epsilon:         gasket.Epsilon,
		}
		return res, flange_model.GasketData_Type(flange_model.GasketData_Type_value[gasket.Type]), nil
	}

	res := &flange_model.GasketResult{
		Gasket:          data.Data.Title,
		Type:            data.Data.Type.String(),
		Thickness:       data.Thickness,
		DOut:            data.DOut,
		Width:           bp,
		M:               data.Data.M,
		Pres:            data.Data.Qo,
		PermissiblePres: data.Data.PermissiblePres,
		Compression:     data.Data.Compression,
		Epsilon:         data.Data.Epsilon,
	}
	return res, data.Data.Type, nil
}
