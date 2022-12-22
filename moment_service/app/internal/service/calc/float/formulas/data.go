package formulas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

func (s *FormulasService) getDataFormulas(data models.DataFloat, result calc_api.FloatResponse, formulas *float_model.Formulas) *float_model.Formulas {
	fWidth := strings.ReplaceAll(strconv.FormatFloat(result.Flange.Width, 'G', 3, 64), "E", "*10^")
	DIn := strings.ReplaceAll(strconv.FormatFloat(result.Flange.DIn, 'G', 3, 64), "E", "*10^")
	gWidth := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Width, 'G', 3, 64), "E", "*10^")
	gDOut := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.DOut, 'G', 3, 64), "E", "*10^")
	b0 := strings.ReplaceAll(strconv.FormatFloat(data.B0, 'G', 3, 64), "E", "*10^")

	if result.Data.HasThorn {
		formulas.B0 = fmt.Sprintf("(%s + %s) / 2", fWidth, gWidth)
		formulas.Dcp = fmt.Sprintf("%s + %s", DIn, fWidth)
	} else {
		if data.TypeGasket != "Soft" {
			formulas.B0 = fmt.Sprintf("%s/4", gWidth)
			formulas.Dcp = fmt.Sprintf("%s - %s/2", gDOut, gWidth)
		} else {
			if data.Gasket.Width > constants.Bp {
				formulas.B0 = fmt.Sprintf("%.1f * sqrt(%s)", constants.B0, gWidth)
			}
			formulas.Dcp = fmt.Sprintf("%s - %s", gDOut, b0)
		}
	}

	return formulas
}
