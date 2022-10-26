package formulas

import (
	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_circle_model"
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

	// занчение зависит от поля "Условия работы"
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

func (s *FormulasService) GetFormulas(req calc_api.ExpressCircleRequest, d models.DataExCircle, result calc_api.ExpressCircleResponse,
) *express_circle_model.Formulas {
	formulas := &express_circle_model.Formulas{
		Deformation:   s.getDeformationFormulas(req, d, result),
		ForsesInBolts: s.getForcesFormulas(req, d, result),
		Bolt:          s.getBoltFormulas(req, d, result),
		Moment:        s.getMomentFormulas(req, d, result),
	}

	return formulas
}
