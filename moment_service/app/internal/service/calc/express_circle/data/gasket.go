package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_circle_model"
)

func (s *DataService) getGasketData(ctx context.Context, gasket *express_circle_model.GasketData,
) (gasketData *express_circle_model.GasketResult, typeGasket express_circle_model.GasketData_Type, err error) {
	width := (gasket.DOut - gasket.DIn) / 2

	if gasket.GasketId != "another" {
		dto := models.GetGasket{GasketId: gasket.GasketId, EnvId: gasket.EnvId, Thickness: gasket.Thickness}
		g, err := s.gasket.GetFullData(ctx, dto)
		if err != nil {
			return nil, 0, err
		}

		res := &express_circle_model.GasketResult{
			Gasket:          g.Gasket,
			Env:             g.Env,
			Type:            g.Type,
			Thickness:       gasket.Thickness,
			DOut:            gasket.DOut,
			DIn:             gasket.DIn,
			Width:           width,
			M:               g.M,
			Pres:            g.SpecificPres,
			PermissiblePres: g.PermissiblePres,
			Compression:     g.Compression,
			Epsilon:         g.Epsilon,
		}
		return res, express_circle_model.GasketData_Type(express_circle_model.GasketData_Type_value[g.Type]), nil
	}

	res := &express_circle_model.GasketResult{
		Gasket:          gasket.Data.Title,
		Type:            gasket.Data.Type.String(),
		Thickness:       gasket.Thickness,
		DIn:             gasket.DIn,
		DOut:            gasket.DOut,
		Width:           width,
		M:               gasket.Data.M,
		Pres:            gasket.Data.Qo,
		PermissiblePres: gasket.Data.PermissiblePres,
		Compression:     gasket.Data.Compression,
		Epsilon:         gasket.Data.Epsilon,
	}
	return res, gasket.Data.Type, nil
}
