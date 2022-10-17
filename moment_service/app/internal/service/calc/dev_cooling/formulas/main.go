package formulas

import (
	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

type FormulasService struct {
	typeBolt map[string]float64
}

func NewFormulasService() *FormulasService {
	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}

	return &FormulasService{
		typeBolt: bolt,
	}
}

func (s *FormulasService) GetFormulas(
	Ab, Lambda1, Lambda2, Alpha1, Alpha2 float64,
	req calc_api.DevCoolingRequest,
	d models.DataDevCooling,
	result calc_api.DevCoolingResponse,
) *dev_cooling_model.Formulas {
	formulas := &dev_cooling_model.Formulas{
		Auxiliary: s.GetAuxiliaryFormulas(req, d, result),
		Bolt:      s.GetBoltFormulas(Lambda1, Lambda2, Alpha1, Alpha2, req, d, result),
	}

	return formulas
}
