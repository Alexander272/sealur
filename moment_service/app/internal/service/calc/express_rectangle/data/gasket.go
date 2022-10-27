package data

import (
	"context"

	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_rectangle_model"
)

func (s *DataService) getGasketData(ctx context.Context, gasket *express_rectangle_model.GasketData,
) (gasketData *express_rectangle_model.GasketResult, typeGasket express_rectangle_model.GasketData_Type, err error) {
	if gasket.GasketId != "another" {
		dto := models.GetGasket{GasketId: gasket.GasketId, EnvId: gasket.EnvId, Thickness: gasket.Thickness}
		g, err := s.gasket.GetFullData(ctx, dto)
		if err != nil {
			return nil, 0, err
		}

		res := &express_rectangle_model.GasketResult{
			Gasket:          g.Gasket,
			Env:             g.Env,
			Type:            g.Type,
			Thickness:       gasket.Thickness,
			SizeLong:        gasket.SizeLong,
			SizeTrans:       gasket.SizeTrans,
			Width:           gasket.Width,
			M:               g.M,
			Pres:            g.SpecificPres,
			PermissiblePres: g.PermissiblePres,
			Compression:     g.Compression,
			Epsilon:         g.Epsilon,
		}
		return res, express_rectangle_model.GasketData_Type(express_rectangle_model.GasketData_Type_value[g.Type]), nil
	}

	res := &express_rectangle_model.GasketResult{
		Gasket:          gasket.Data.Title,
		Type:            gasket.Data.Type.String(),
		Thickness:       gasket.Thickness,
		SizeLong:        gasket.SizeLong,
		SizeTrans:       gasket.SizeTrans,
		Width:           gasket.Width,
		M:               gasket.Data.M,
		Pres:            gasket.Data.Qo,
		PermissiblePres: gasket.Data.PermissiblePres,
		Compression:     gasket.Data.Compression,
		Epsilon:         gasket.Data.Epsilon,
	}
	return res, gasket.Data.Type, nil
}
