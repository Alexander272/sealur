package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

func (s *DataService) getGasketData(ctx context.Context, data *cap_model.GasketData, bp float64) (*cap_model.GasketResult, cap_model.GasketData_Type, error) {
	if data.GasketId != "another" {
		g := models.GetGasket{GasketId: data.GasketId, EnvId: data.EnvId, Thickness: data.Thickness}
		gasket, err := s.gasket.GetFullData(ctx, g)
		if err != nil {
			return nil, 0, err
		}

		res := &cap_model.GasketResult{
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
		return res, cap_model.GasketData_Type(cap_model.GasketData_Type_value[gasket.Type]), nil
	}

	res := &cap_model.GasketResult{
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
