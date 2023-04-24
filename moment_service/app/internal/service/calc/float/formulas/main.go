package formulas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
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

func (s *FormulasService) GetFormulas(
	req *calc_api.FloatRequest,
	Condition, TypeBolt string,
	IsWork bool,
	data models.DataFloat,
	result calc_api.FloatResponse,
) *float_model.Formulas {
	formulas := &float_model.Formulas{}

	friction := strconv.FormatFloat(req.Friction, 'G', 3, 64)
	area := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Area, 'G', 3, 64), "E", "*10^")
	diameter := strconv.FormatFloat(data.Bolt.Diameter, 'G', 3, 64)
	Lb0 := strings.ReplaceAll(strconv.FormatFloat(result.Bolt.Lenght, 'G', 3, 64), "E", "*10^")
	typeBolt := strings.ReplaceAll(strconv.FormatFloat(s.typeBolt[TypeBolt], 'G', 3, 64), "E", "*10^")
	bEpsilonAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.EpsilonAt20, 'G', 3, 64), "E", "*10^")
	bSigmaAt20 := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.SigmaAt20, 'G', 3, 64), "E", "*10^")
	bSigma := strings.ReplaceAll(strconv.FormatFloat(data.Bolt.Sigma, 'G', 3, 64), "E", "*10^")
	count := data.Bolt.Count

	compression := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Compression, 'G', 3, 64), "E", "*10^")
	gEpsilon := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Epsilon, 'G', 3, 64), "E", "*10^")
	gThickness := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Thickness, 'G', 3, 64), "E", "*10^")
	gWidth := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Width, 'G', 3, 64), "E", "*10^")
	gM := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.M, 'G', 3, 64), "E", "*10^")
	gPres := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.Pres, 'G', 3, 64), "E", "*10^")
	gPermissiblePres := strings.ReplaceAll(strconv.FormatFloat(data.Gasket.PermissiblePres, 'G', 3, 64), "E", "*10^")

	fB := strings.ReplaceAll(strconv.FormatFloat(data.Flange.B, 'G', 3, 64), "E", "*10^")

	cY := strings.ReplaceAll(strconv.FormatFloat(data.Cap.Y, 'G', 3, 64), "E", "*10^")

	Dcp := strings.ReplaceAll(strconv.FormatFloat(data.Dcp, 'G', 3, 64), "E", "*10^")
	b0 := strings.ReplaceAll(strconv.FormatFloat(data.B0, 'G', 3, 64), "E", "*10^")
	pressure := strings.ReplaceAll(strconv.FormatFloat(math.Abs(result.Data.Pressure), 'G', 3, 64), "E", "*10^")
	Lb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Lb, 'G', 3, 64), "E", "*10^")
	yp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Yp, 'G', 3, 64), "E", "*10^")
	yb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Yb, 'G', 3, 64), "E", "*10^")
	Po := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Po, 'G', 3, 64), "E", "*10^")
	A := strings.ReplaceAll(strconv.FormatFloat(result.Calc.A, 'G', 3, 64), "E", "*10^")
	alpha := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Alpha, 'G', 3, 64), "E", "*10^")
	Qd := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Qd, 'G', 3, 64), "E", "*10^")
	Rp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Rp, 'G', 3, 64), "E", "*10^")
	Pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Pb1, 'G', 3, 64), "E", "*10^")
	Pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Pb2, 'G', 3, 64), "E", "*10^")
	Pb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Pb, 'G', 3, 64), "E", "*10^")
	Pbr := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Pbr, 'G', 3, 64), "E", "*10^")
	dSigmaM := strings.ReplaceAll(strconv.FormatFloat(result.Calc.DSigmaM, 'G', 3, 64), "E", "*10^")
	Mkp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Mkp, 'G', 3, 64), "E", "*10^")

	formulas = s.getDataFormulas(data, result, formulas)
	formulas.Flange = s.getFlangeFormulas(result.Flange, result.Cap, Dcp)
	formulas.Cap = s.getCapFormulas(result.Flange, result.Cap)

	if data.TypeGasket == "Soft" {
		formulas.Yp = fmt.Sprintf("(%s * %s) / (%s * %f * %s *%s)", gThickness, compression, gEpsilon, math.Pi, Dcp, gWidth)
	}
	formulas.Lb = fmt.Sprintf("%s + %s * %s", Lb0, typeBolt, diameter)
	formulas.Yb = fmt.Sprintf("%s / (%s * %s * %d)", Lb, bEpsilonAt20, area, count)

	formulas.A = fmt.Sprintf("%d * %s", count, area)
	formulas.Po = fmt.Sprintf("0.5 * %f * %s * %s * %s", math.Pi, Dcp, b0, gPres)
	if result.Data.Pressure >= 0 {
		formulas.Rp = fmt.Sprintf("%f * %s * %s * %s * |%s|", math.Pi, Dcp, b0, gM, pressure)
	}
	formulas.Qd = fmt.Sprintf("0.785 * (%s)^2 * %s", Dcp, pressure)
	if data.TypeGasket == "Oval" {
		formulas.Alpha = fmt.Sprintf("1 - (%s - (%s - %s)*%s)/(%s + %s)", yp, cY, fB, fB, yp, yb)
	}

	formulas.Pb2 = fmt.Sprintf("max(%s; 0.4 * %s * %s)", Po, A, bSigmaAt20)
	formulas.Pb1 = fmt.Sprintf("%s*%s + %s", alpha, Qd, Rp)

	formulas.Pb = fmt.Sprintf("max(%s; %s)", Pb1, Pb2)
	formulas.Pbr = fmt.Sprintf("%s + (1 - %s)*%s", Pb, alpha, Qd)

	formulas.SigmaB1 = fmt.Sprintf("%s / %s", Pb, A)
	formulas.SigmaB2 = fmt.Sprintf("%s / %s", Pbr, A)

	Kyp := strconv.FormatFloat(s.Kyp[IsWork], 'G', 3, 64)
	Kyz := strconv.FormatFloat(s.Kyz[Condition], 'G', 3, 64)
	Kyt := strconv.FormatFloat(constants.NoLoadKyt, 'G', 3, 64)

	formulas.DSigmaM = fmt.Sprintf("1.2 * %s * %s * %s * %s", Kyp, Kyz, Kyt, bSigmaAt20)
	formulas.DSigmaR = fmt.Sprintf("%s * %s * %s * %s", Kyp, Kyz, Kyt, bSigma)

	if data.TypeGasket == "Soft" {
		formulas.Qmax = fmt.Sprintf("max(%s; %s) / (%f * %s * %s)", Pb, Pbr, math.Pi, Dcp, gWidth)
	}

	if !(result.Calc.SigmaB1 > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter && data.Bolt.Diameter <= constants.MaxDiameter) {
		formulas.Mkp = fmt.Sprintf("(%s * %s * %s/%d) / 1000", friction, Pb, diameter, count)
	}
	formulas.Mkp1 = fmt.Sprintf("0.75 * %s", Mkp)

	Prek := fmt.Sprintf("0.8 * %s * %s", A, bSigmaAt20)
	formulas.Qrek = fmt.Sprintf("(%s) / (%f * %s * %s)", Prek, math.Pi, Dcp, gWidth)
	formulas.Mrek = fmt.Sprintf("(%s * %s * %s/%d) / 1000", friction, Prek, diameter, count)

	Pmax := fmt.Sprintf("%s * %s", dSigmaM, A)
	formulas.Qmax = fmt.Sprintf("(%s) / (%f * %s *%s)", Pmax, math.Pi, Dcp, gWidth)

	if data.TypeGasket == "Soft" && result.Calc.Qmax > data.Gasket.PermissiblePres {
		Pmax = fmt.Sprintf("%s * %f * %s *%s", gPermissiblePres, math.Pi, Dcp, gWidth)
	}

	if result.Calc.Mrek > result.Calc.Mmax {
		formulas.Mrek = ""
	}
	if result.Calc.Qrek > result.Calc.Qmax {
		formulas.Qrek = ""
	}

	formulas.Mmax = fmt.Sprintf("(%s * %s *%s / %d) / 1000", friction, Pmax, diameter, count)

	return formulas
}
