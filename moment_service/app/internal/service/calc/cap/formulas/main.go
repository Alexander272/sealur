package formulas

import (
	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

type FormulasService struct {
	typeBolt map[string]float64
	Kyp      map[bool]float64
	Kyz      map[string]float64
	Kyt      map[bool]float64
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

	kt := map[bool]float64{
		true:  constants.LoadKyt,
		false: constants.NoLoadKyt,
	}

	return &FormulasService{
		typeBolt: bolt,
		Kyp:      kp,
		Kyz:      kz,
		Kyt:      kt,
	}
}

func (s *FormulasService) GetFormulas(req *calc_api.CapRequest, d models.DataCap, result *calc_api.CapResponse, aux *cap_model.CalcAuxiliary,
) *cap_model.Formulas {
	formulas := &cap_model.Formulas{}

	if req.Data.Calculation == cap_model.MainData_basis {
		formulas.Basis = s.basisFormulas(req, d, result, aux)
	} else {
		formulas.Strength = s.strengthFormulas(req, d, result)
	}

	return formulas
}
