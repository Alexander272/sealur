package formulas

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

func (s *FormulasService) basisFormulas(
	req *calc_api.CapRequest,
	d models.DataCap,
	result *calc_api.CapResponse,
	aux *cap_model.CalcAuxiliary,
) *cap_model.Formulas_Basis {
	bolt := s.boltStrengthFormulas(
		req, d, result,
		result.Calc.Basis.ForcesInBolts.Pb,
		result.Calc.Basis.ForcesInBolts.Pbr,
		result.Calc.Basis.ForcesInBolts.A,
		result.Calc.Basis.Deformation.Dcp,
		true,
	)

	moment := s.momentFormulas(
		req, d, result,
		result.Calc.Basis.BoltStrength.SigmaB1,
		result.Calc.Basis.BoltStrength.DSigmaM,
		result.Calc.Basis.ForcesInBolts.Pb,
		result.Calc.Basis.ForcesInBolts.A,
		result.Calc.Basis.Deformation.Dcp,
		result.Calc.Basis.Moment,
		true,
	)

	formulas := &cap_model.Formulas_Basis{
		Deformation:   s.deformationFormulas(req, d, result),
		ForcesInBolts: s.forcesInBoltsCalculate(req, d, result, aux),
		BoltStrength:  bolt,
		Moment:        moment,
	}

	return formulas
}

func (s *FormulasService) deformationFormulas(req *calc_api.CapRequest, d models.DataCap, result *calc_api.CapResponse,
) *cap_model.DeformationFormulas {
	deformation := &cap_model.DeformationFormulas{}

	// перевод чисел в строки
	pressure := strconv.FormatFloat(req.Data.Pressure, 'f', -1, 64)

	width := strconv.FormatFloat(d.Gasket.Width, 'f', -1, 64)
	dOut := strconv.FormatFloat(d.Gasket.DOut, 'f', -1, 64)
	pres := strconv.FormatFloat(d.Gasket.Pres, 'f', -1, 64)
	m := strconv.FormatFloat(d.Gasket.M, 'f', -1, 64)

	B0 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.B0, 'G', 3, 64), "E", "*10^")
	Dcp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.Dcp, 'G', 3, 64), "E", "*10^")

	if d.TypeGasket == cap_model.GasketData_Oval {
		// формула 4
		deformation.B0 = fmt.Sprintf("%s / 4", width)
		// формула ?
		deformation.Dcp = fmt.Sprintf("%s - %s/2", dOut, width)

	} else {
		if !(d.Gasket.Width <= constants.Bp) {
			// формула 3
			deformation.B0 = fmt.Sprintf("%.1f * sqrt(%s)", constants.B0, width)
		}
		// формула 5
		deformation.Dcp = fmt.Sprintf("%s - %s", dOut, B0)
	}

	// формула 6
	// Усилие необходимое для смятия прокладки при затяжке
	deformation.Po = fmt.Sprintf("0.5 * %f * %s * %s * %s", math.Pi, Dcp, B0, pres)

	if req.Data.Pressure >= 0 {
		// формула 7
		// Усилие на прокладке в рабочих условиях
		deformation.Rp = fmt.Sprintf("%f * %s * %s * %s * |%s|", math.Pi, Dcp, B0, m, pressure)
	}

	return deformation
}

func (s *FormulasService) forcesInBoltsCalculate(
	req *calc_api.CapRequest, d models.DataCap, result *calc_api.CapResponse, aux *cap_model.CalcAuxiliary,
) *cap_model.ForcesInBoltsFormulas {
	forces := &cap_model.ForcesInBoltsFormulas{}

	// перевод чисел в строки
	axialForce := req.Data.AxialForce
	pressure := strconv.FormatFloat(req.Data.Pressure, 'f', -1, 64)
	temp := strconv.FormatFloat(req.Data.Temp, 'f', -1, 64)

	count := d.Bolt.Count
	area := strconv.FormatFloat(d.Bolt.Area, 'G', 3, 64)
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 3, 64)
	bAlpha := strings.ReplaceAll(strconv.FormatFloat(d.Bolt.Alpha, 'G', 3, 64), "E", "*10^")
	bTemp := strconv.FormatFloat(d.Bolt.Temp, 'G', 3, 64)

	alphaF1 := strings.ReplaceAll(strconv.FormatFloat(d.Flange.Alpha, 'G', 3, 64), "E", "*10^")
	alphaF2 := strings.ReplaceAll(strconv.FormatFloat(d.Cap.Alpha, 'G', 3, 64), "E", "*10^")
	h1 := strconv.FormatFloat(d.Flange.H, 'G', 3, 64)
	h2 := strconv.FormatFloat(d.Cap.H, 'G', 3, 64)
	tf1 := strconv.FormatFloat(d.Flange.T, 'G', 3, 64)
	tf2 := strconv.FormatFloat(d.Cap.T, 'G', 3, 64)

	var wAlpha1, wAlpha2, thick1, thick2 string
	if req.IsUseWasher {
		wAlpha1 = strings.ReplaceAll(strconv.FormatFloat(d.Washer1.Alpha, 'G', 3, 64), "E", "*10^")
		wAlpha2 = strings.ReplaceAll(strconv.FormatFloat(d.Washer2.Alpha, 'G', 3, 64), "E", "*10^")
		thick1 = strconv.FormatFloat(d.Washer1.Thickness, 'G', 3, 64)
		thick2 = strconv.FormatFloat(d.Washer2.Thickness, 'G', 3, 64)
	}

	var eAlpha, eThick string
	if req.Data.IsEmbedded {
		eAlpha = strings.ReplaceAll(strconv.FormatFloat(d.Embed.Alpha, 'G', 3, 64), "E", "*10^")
		eThick = strconv.FormatFloat(d.Embed.Thickness, 'G', 3, 64)
	}

	Dcp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.Dcp, 'G', 3, 64), "E", "*10^")
	Po := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.Po, 'G', 3, 64), "E", "*10^")
	Rp := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.Deformation.Rp, 'G', 3, 64), "E", "*10^")

	Qd := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Qd, 'G', 3, 64), "E", "*10^")
	Ab := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.A, 'G', 3, 64), "E", "*10^")
	Alpha := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Alpha, 'G', 3, 64), "E", "*10^")
	Qt := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Qt, 'G', 3, 64), "E", "*10^")
	Pb1 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Pb1, 'G', 3, 64), "E", "*10^")
	Pb2 := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Pb2, 'G', 3, 64), "E", "*10^")
	Pb := strings.ReplaceAll(strconv.FormatFloat(result.Calc.Basis.ForcesInBolts.Pb, 'G', 3, 64), "E", "*10^")

	Yp := strings.ReplaceAll(strconv.FormatFloat(aux.Yp, 'G', 3, 64), "E", "*10^")
	Yb := strings.ReplaceAll(strconv.FormatFloat(aux.Yb, 'G', 3, 64), "E", "*10^")
	Yf1 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange.Yf, 'G', 3, 64), "E", "*10^")
	E1 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange.E, 'G', 3, 64), "E", "*10^")
	B1 := strings.ReplaceAll(strconv.FormatFloat(aux.Flange.B, 'G', 3, 64), "E", "*10^")
	Y := strings.ReplaceAll(strconv.FormatFloat(aux.Cap.Y, 'G', 3, 64), "E", "*10^")
	Gamma := strings.ReplaceAll(strconv.FormatFloat(aux.Gamma, 'G', 3, 64), "E", "*10^")

	// формула 8
	// Суммарная площадь сечения болтов/шпилек
	forces.A = fmt.Sprintf("%d * %s", count, area)

	// формула 9
	// Равнодействующая нагрузка от давления
	forces.Qd = fmt.Sprintf("0.785 * (%s)^2 * %s", Dcp, pressure)

	// формула 10
	// Приведенная нагрузка, вызванная воздействием внешней силы и изгибающего момента
	// forces.Qfm = fmt.Sprintf("%d", axialForce)

	if !(d.TypeGasket == cap_model.GasketData_Oval || d.FlangeType == cap_model.FlangeData_free) {
		// формула (Е.11)
		// Коэффициент жесткости
		forces.Alpha = fmt.Sprintf("1 - (%s - (%s * %s + %s * %s) * %s)/(%s + %s + (%s + %s) * (%s)^2)",
			Yp, Yf1, E1, Y, B1, B1, Yp, Yb, Yf1, Y, B1)
	}

	minB := fmt.Sprintf("0.4 * %s * %s", Ab, sigmaAt20)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	forces.Pb2 = fmt.Sprintf("max(%s; %s)", Po, minB)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения в рабочих условиях давления на
	// прокладку достаточного для герметизации фланцевого соединения
	forces.Pb1 = fmt.Sprintf("%s * (%s + %d) + %s", Alpha, Qd, axialForce, Rp)

	var temp1, temp2 string
	if req.IsUseWasher {
		temp1 = fmt.Sprintf("(%s * %s + %s * %s) * (%s - 20) + (%s * %s + %s * %s) * (%s - 20)",
			alphaF1, h1, wAlpha1, thick1, tf1, alphaF2, h2, wAlpha2, thick2, tf2)
	} else {
		temp1 = fmt.Sprintf("%s * %s * (%s - 20) + %s * %s *(%s - 20)", alphaF1, h1, tf1, alphaF2, h2, tf2)
	}
	temp2 = fmt.Sprintf("%s + %s", h1, h2)

	if d.FlangeType == cap_model.FlangeData_free {
		alphaK1 := strings.ReplaceAll(strconv.FormatFloat(d.Flange.Ring.Alpha, 'G', 3, 64), "E", "*10^")
		hk1 := strconv.FormatFloat(d.Flange.Ring.Hk, 'G', 3, 64)
		tk1 := strconv.FormatFloat(d.Flange.Ring.T, 'G', 3, 64)

		temp1 += fmt.Sprintf(" + (%s * %s) * (%s - 20)", alphaK1, hk1, tk1)
		temp2 += fmt.Sprintf(" + %s", hk1)
	}

	if req.Data.IsEmbedded {
		temp1 += fmt.Sprintf(" + (%s * %s) * (%s - 20)", eAlpha, eThick, temp)
		temp2 += fmt.Sprintf(" + %s", eThick)
	}

	//? должно быть два варианта формулы с шайбой и без нее
	// шайба будет задаваться так же как и болты + толщина шайбы

	//формула 11 (в старом 13)
	forces.Qt = fmt.Sprintf("%s * ((%s) - %s * (%s) * (%s - 20))", Gamma, temp1, bAlpha, temp2, bTemp)

	forces.Pb1 = fmt.Sprintf("max(%s; %s-%s)", forces.Pb1, forces.Pb1, Qt)
	forces.Pb = fmt.Sprintf("max(%s; %s)", Pb1, Pb2)
	forces.Pbr = fmt.Sprintf("%s + (1 - %s) * (%s + %d) + %s", Pb, Alpha, Qd, axialForce, Qt)

	return forces
}

func (s *FormulasService) boltStrengthFormulas(
	req *calc_api.CapRequest, d models.DataCap, result *calc_api.CapResponse,
	Pbm, Pbr, Ab, Dcp float64,
	isLoad bool,
) *cap_model.BoltStrengthFormulas {
	bolt := &cap_model.BoltStrengthFormulas{}

	// перевод чисел в строки
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 3, 64)
	sigma := strconv.FormatFloat(d.Bolt.Sigma, 'G', 3, 64)

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 3, 64)

	Ab_ := strings.ReplaceAll(strconv.FormatFloat(Ab, 'G', 3, 64), "E", "*10^")
	Pb_ := strings.ReplaceAll(strconv.FormatFloat(Pbm, 'G', 3, 64), "E", "*10^")
	Pbr_ := strings.ReplaceAll(strconv.FormatFloat(Pbr, 'G', 3, 64), "E", "*10^")

	Dcp_ := strings.ReplaceAll(strconv.FormatFloat(Dcp, 'G', 3, 64), "E", "*10^")

	bolt.SigmaB1 = fmt.Sprintf("%s / %s", Pb_, Ab_)
	bolt.SigmaB2 = fmt.Sprintf("%s / %s", Pbr_, Ab_)

	Kyp := strconv.FormatFloat(s.Kyp[req.Data.IsWork], 'G', 3, 64)
	Kyz := strconv.FormatFloat(s.Kyz[req.Data.Condition.String()], 'G', 3, 64)
	Kyt := strconv.FormatFloat(s.Kyt[isLoad], 'G', 3, 64)

	// формула Г.3
	bolt.DSigmaM = fmt.Sprintf("1.2 * %s * %s * %s * %s", Kyp, Kyz, Kyt, sigmaAt20)
	// формула Г.4
	bolt.DSigmaR = fmt.Sprintf("%s * %s * %s * %s", Kyp, Kyz, Kyt, sigma)

	if d.TypeGasket == cap_model.GasketData_Soft {
		bolt.Q = fmt.Sprintf("max(%s; %s) / (%f * %s * %s)", Pb_, Pbr_, math.Pi, Dcp_, width)
	}

	return bolt
}

func (s *FormulasService) momentFormulas(
	req *calc_api.CapRequest, d models.DataCap, result *calc_api.CapResponse,
	SigmaB1, DSigmaM, Pbm, Ab, Dcp float64,
	mom *cap_model.CalcMoment,
	fullCalculate bool,
) *cap_model.MomentFormulas {
	moment := &cap_model.MomentFormulas{}

	// перевод чисел в строки
	sigmaAt20 := strconv.FormatFloat(d.Bolt.SigmaAt20, 'G', 3, 64)
	diameter := strconv.FormatFloat(d.Bolt.Diameter, 'G', 3, 64)
	count := d.Bolt.Count

	width := strconv.FormatFloat(d.Gasket.Width, 'G', 3, 64)
	perPres := strconv.FormatFloat(d.Gasket.PermissiblePres, 'G', 3, 64)

	Ab_ := strings.ReplaceAll(strconv.FormatFloat(Ab, 'G', 3, 64), "E", "*10^")
	Pb_ := strings.ReplaceAll(strconv.FormatFloat(Pbm, 'G', 3, 64), "E", "*10^")
	Dcp_ := strings.ReplaceAll(strconv.FormatFloat(Dcp, 'G', 3, 64), "E", "*10^")
	DSigmaM_ := strings.ReplaceAll(strconv.FormatFloat(DSigmaM, 'G', 3, 64), "E", "*10^")
	Mkp := strings.ReplaceAll(strconv.FormatFloat(mom.Mkp, 'G', 3, 64), "E", "*10^")

	if !(SigmaB1 > constants.MaxSigmaB && d.Bolt.Diameter >= constants.MinDiameter && d.Bolt.Diameter <= constants.MaxDiameter) {
		moment.Mkp = fmt.Sprintf("(0.3 * %s * %s / %d) / 1000", Pb_, diameter, count)
	}

	moment.Mkp1 = fmt.Sprintf("0.75 * %s", Mkp)

	if fullCalculate {
		Prek := fmt.Sprintf("0.8 * %s * %s", Ab_, sigmaAt20)
		moment.Qrek = fmt.Sprintf("%s / (%f * %s * %s)", Prek, math.Pi, Dcp_, width)
		moment.Mrek = fmt.Sprintf("(0.3 * %s * %s / %d) / 1000", Prek, diameter, count)

		Pmax := fmt.Sprintf("%s * %s", DSigmaM_, Ab_)
		moment.Qmax = fmt.Sprintf("%s / (%f * %s * %s)", Pmax, math.Pi, Dcp_, width)

		if d.TypeGasket == cap_model.GasketData_Soft && mom.Qmax > d.Gasket.PermissiblePres {
			Pmax = fmt.Sprintf("%s * (%f * %s * %s)", perPres, math.Pi, Dcp_, width)
		}

		moment.Mmax = fmt.Sprintf("(0.3 * %s * %s / %d) / 1000", Pmax, diameter, count)
	}

	return moment
}
