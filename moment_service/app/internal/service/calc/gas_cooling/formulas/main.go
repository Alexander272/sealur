package formulas

import (
	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/gas_cooling_model"
)

type FormulasService struct {
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
}

func NewFormulasService() *FormulasService {
	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}

	// значение зависит от поля "Условия работы"
	kp := map[bool]float64{
		true:  constants.WorkKyp,
		false: constants.TestKyp,
	}

	// значение зависит от поля "Условие затяжки"
	kz := map[string]float64{
		"uncontrollable":  constants.UncontrollableKyz,
		"controllable":    constants.ControllableKyz,
		"controllablePin": constants.ControllablePinKyz,
	}

	return &FormulasService{
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
	}
}

func (s *FormulasService) GetFormulas(req *calc_api.GasCoolingRequest, d models.DataGasCooling, result *calc_api.GasCoolingResponse,
) *gas_cooling_model.Formulas {
	formulas := &gas_cooling_model.Formulas{
		Auxiliary:     s.getAuxiliaryFormulas(req, d, result),
		ForcesInBolts: s.getForcesFormulas(req, d, result),
		Bolt:          s.getBoltFormulas(req, d, result),
		Moment:        s.getMomentFormulas(req, d, result),
	}

	return formulas
}
