package formulas

import (
	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
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
	kp := map[bool]float64{
		true:  constants.WorkKyp,
		false: constants.TestKyp,
	}
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

func (s *FormulasService) GetFormulas(req *calc_api.FlangeRequest, d models.DataFlange, result *calc_api.FlangeResponse, aux *flange_model.CalcAuxiliary,
) *flange_model.Formulas {
	formulas := &flange_model.Formulas{}

	if req.Calculation == calc_api.FlangeRequest_basis {
		formulas.Basis = s.basisFormulas(req, d, result, aux)
	} else {
		formulas.Strength = s.strengthFormulas(req, d, result)
	}

	return formulas
}
