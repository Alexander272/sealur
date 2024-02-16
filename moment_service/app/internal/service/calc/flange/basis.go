package flange

import (
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
)

// Расчет основных величин
func (s *FlangeService) basisCalculate(data models.DataFlange, req *calc_api.FlangeRequest) (*flange_model.Calculated_Basis, *flange_model.CalcAuxiliary) {
	deformation := s.deformationCalculate(data, req.Pressure)
	forces, aux := s.forcesInBoltsCalculate(data, deformation, req)
	bolts := s.boltStrengthCalculate(data, req, forces.Pb, forces.Pbr, forces.A, deformation.Dcp, true)
	moment := &flange_model.CalcMoment{}

	ok := (bolts.VSigmaB1 && bolts.VSigmaB2 && data.TypeGasket != flange_model.GasketData_Soft) ||
		(bolts.VSigmaB1 && bolts.VSigmaB2 && bolts.Q <= float64(data.Gasket.PermissiblePres) && data.TypeGasket == flange_model.GasketData_Soft)
	if ok {
		moment = s.momentCalculate(req.Friction, data, bolts.SigmaB1, bolts.DSigmaM, forces.Pb, forces.A, deformation.Dcp, true)
	}

	res := &flange_model.Calculated_Basis{
		Deformation:   deformation,
		ForcesInBolts: forces,
		BoltStrength:  bolts,
		Moment:        moment,
	}

	return res, aux
}

// Усилия, необходимые для смятия прокладки и обеспечения герметичности фланцевого соединения
func (s *FlangeService) deformationCalculate(data models.DataFlange, pressure float64) *flange_model.CalcDeformation {
	deformation := &flange_model.CalcDeformation{}

	if data.TypeGasket == flange_model.GasketData_Oval {
		// формула 4
		deformation.B0 = data.Gasket.Width / 4
		// формула ?
		deformation.Dcp = data.Gasket.DOut - data.Gasket.Width/2

	} else {
		if data.Gasket.Width <= constants.Bp {
			// формула 2
			deformation.B0 = data.Gasket.Width
		} else {
			// формула 3
			deformation.B0 = constants.B0 * math.Sqrt(data.Gasket.Width)
		}
		// формула 5
		deformation.Dcp = data.Gasket.DOut - deformation.B0
	}

	// формула 6
	// Усилие необходимое для смятия прокладки при затяжке
	deformation.Po = 0.5 * math.Pi * deformation.Dcp * deformation.B0 * data.Gasket.Pres

	if pressure >= 0 {
		// формула 7
		// Усилие на прокладке в рабочих условиях
		deformation.Rp = math.Pi * deformation.Dcp * deformation.B0 * data.Gasket.M * math.Abs(pressure)
	}

	return deformation
}

// Усилия в болтах (шпильках) фланцевого соединения при затяжке и в рабочих условиях
func (s *FlangeService) forcesInBoltsCalculate(
	data models.DataFlange,
	def *flange_model.CalcDeformation,
	req *calc_api.FlangeRequest,
) (*flange_model.CalcForcesInBolts, *flange_model.CalcAuxiliary) {
	forces := &flange_model.CalcForcesInBolts{}

	// формула 8
	// Суммарная площадь сечения болтов/шпилек
	forces.A = float64(data.Bolt.Count) * data.Bolt.Area

	// формула 9
	// Равнодействующая нагрузка от давления
	forces.Qd = 0.785 * math.Pow(def.Dcp, 2) * req.Pressure

	temp1 := float64(req.AxialForce) + 4*math.Abs(float64(req.BendingMoment))/def.Dcp
	temp2 := float64(req.AxialForce) - 4*math.Abs(float64(req.BendingMoment))/def.Dcp

	// формула 10
	// Приведенная нагрузка, вызванная воздействием внешней силы и изгибающего момента
	forces.Qfm = math.Max(temp1, temp2)

	var yp float64 = 0
	if data.TypeGasket == flange_model.GasketData_Soft {
		// Податливость прокладки
		yp = (data.Gasket.Thickness * data.Gasket.Compression) / (data.Gasket.Epsilon * math.Pi * def.Dcp * data.Gasket.Width)
	}

	// приложение К пояснение к формуле К.2
	Lb := data.Bolt.Length + s.typeBolt[req.Type.String()]*data.Bolt.Diameter
	// формула К.2
	// Податливость болтов/шпилек
	yb := Lb / (data.Bolt.EpsilonAt20 * data.Bolt.Area * float64(data.Bolt.Count))

	flange1 := s.auxFlangeCalculate(req.FlangesData[0].Type, data.Flange1, def.Dcp)
	flange2 := flange1
	if len(req.FlangesData) > 1 {
		flange2 = s.auxFlangeCalculate(req.FlangesData[1].Type, data.Flange2, def.Dcp)
	}
	aux := &flange_model.CalcAuxiliary{
		Flange1: flange1,
		Flange2: flange2,
	}
	aux.Yp = yp
	aux.Lb = Lb
	aux.Yb = yb

	if data.TypeGasket == flange_model.GasketData_Oval || data.Type1 == flange_model.FlangeData_free || data.Type2 == flange_model.FlangeData_free {
		// Для фланцев с овальными и восьмигранными прокладками и для свободных фланцев коэффициенты жесткости фланцевого соединения принимают равными 1.
		forces.Alpha = 1
	} else {
		// формула (Е.11)
		// Коэффициент жесткости
		forces.Alpha = 1 - (yp-(flange1.Yf*flange1.E*flange1.B+flange2.Yf*flange2.E*flange2.B))/
			(yp+yb+(flange1.Yf*math.Pow(flange1.B, 2)+flange2.Yf*math.Pow(flange2.B, 2)))
	}

	dividend := yb + flange1.Yfn*flange1.B*(flange1.B+flange1.E-math.Pow(flange1.E, 2)/def.Dcp) +
		+flange2.Yfn*flange2.B*(flange2.B+flange2.E-math.Pow(flange2.E, 2)/def.Dcp)
	divider := yb + yp*math.Pow(data.Flange1.D6/def.Dcp, 2) + flange1.Yfn*math.Pow(flange1.B, 2) + flange2.Yfn*math.Pow(flange2.B, 2)

	if data.Type1 == flange_model.FlangeData_free {
		dividend += flange1.Yfc * math.Pow(flange1.A, 2)
		divider += flange1.Yfc * math.Pow(flange1.A, 2)
	}
	if data.Type2 == flange_model.FlangeData_free {
		dividend += flange2.Yfc * math.Pow(flange2.A, 2)
		divider += flange2.Yfc * math.Pow(flange2.A, 2)
	}

	// формула (Е.13)
	// Коэффициент жесткости фланцевого соединения нагруженного внешним изгибающим моментом
	forces.AlphaM = dividend / divider

	minB := 0.4 * forces.A * data.Bolt.SigmaAt20
	forces.MinB = minB
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	forces.Pb2 = math.Max(def.Po, minB)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения в рабочих условиях давления на
	// прокладку достаточного для герметизации фланцевого соединения
	forces.Pb1 = forces.Alpha*(forces.Qd+float64(req.AxialForce)) + def.Rp + 4*forces.AlphaM*math.Abs(float64(req.BendingMoment))/def.Dcp

	divider = yp + yb*data.Bolt.EpsilonAt20/data.Bolt.Epsilon + (flange1.Yf*data.Flange1.EpsilonAt20/data.Flange1.Epsilon)*math.Pow(flange1.B, 2) +
		+(flange2.Yf*data.Flange2.EpsilonAt20/data.Flange2.Epsilon)*math.Pow(flange2.B, 2)

	if data.Type1 == flange_model.FlangeData_free {
		divider += (flange1.Yk * data.Flange1.Ring.EpsilonKAt20 / data.Flange1.Ring.EpsilonK) * math.Pow(flange1.A, 2)
	}
	if data.Type2 == flange_model.FlangeData_free {
		divider += (flange2.Yk * data.Flange2.Ring.EpsilonKAt20 / data.Flange2.Ring.EpsilonK) * math.Pow(flange2.A, 2)
	}

	// формула (Е.8)
	gamma := 1 / divider
	aux.Gamma = gamma

	if req.IsUseWasher {
		temp1 = (data.Flange1.AlphaF*data.Flange1.H+data.Washer1.Alpha*data.Washer1.Thickness)*(data.Flange1.Tf-20) +
			+(data.Flange2.AlphaF*data.Flange2.H+data.Washer2.Alpha*data.Washer2.Thickness)*(data.Flange2.Tf-20)
	} else {
		temp1 = data.Flange1.AlphaF*data.Flange1.H*(data.Flange1.Tf-20) + data.Flange2.AlphaF*data.Flange2.H*(data.Flange2.Tf-20)
	}
	temp2 = data.Flange1.H + data.Flange2.H

	if data.Type1 == flange_model.FlangeData_free {
		temp1 += data.Flange1.Ring.AlphaK * data.Flange1.Hk * (data.Flange1.Ring.Tk - 20)
		temp2 += data.Flange1.Hk
	}
	if data.Type2 == flange_model.FlangeData_free {
		temp1 += data.Flange2.Ring.AlphaK * data.Flange2.Hk * (data.Flange2.Ring.Tk - 20)
		temp2 += data.Flange2.Hk
	}
	if req.IsEmbedded {
		temp1 += data.Embed.Alpha * data.Embed.Thickness * (req.Temp - 20)
		temp2 += data.Embed.Thickness
	}

	//? должно быть два варианта формулы с шайбой и без нее
	// шайба будет задаваться так же как и болты + толщина шайбы

	//формула 11 (в старом 13)
	Qt := gamma * (temp1 - data.Bolt.Alpha*temp2*(data.Bolt.Temp-20))
	forces.Qt = Qt

	forces.Pb1 = math.Max(forces.Pb1, forces.Pb1-Qt)
	forces.Pb = math.Max(forces.Pb1, forces.Pb2)
	forces.Pbr = forces.Pb + (1-forces.Alpha)*(forces.Qd+float64(req.AxialForce)) + Qt + 4*(1-forces.AlphaM*math.Abs(float64(req.BendingMoment)))/def.Dcp

	return forces, aux
}

// Проверка прочности болтов (шпилек) и прокладки
func (s *FlangeService) boltStrengthCalculate(
	data models.DataFlange,
	req *calc_api.FlangeRequest,
	Pbm, Pbr, Ab, Dcp float64,
	isLoad bool,
) *flange_model.CalcBoltStrength {
	bolt := &flange_model.CalcBoltStrength{}

	bolt.SigmaB1 = Pbm / Ab
	bolt.SigmaB2 = Pbr / Ab

	Kyp := s.Kyp[req.IsWork]
	Kyz := s.Kyz[req.Condition.String()]
	Kyt := s.Kyt[isLoad]
	// формула Г.3
	bolt.DSigmaM = 1.2 * Kyp * Kyz * Kyt * data.Bolt.SigmaAt20
	// формула Г.4
	bolt.DSigmaR = Kyp * Kyz * Kyt * data.Bolt.Sigma

	if data.TypeGasket == flange_model.GasketData_Soft {
		bolt.Q = math.Max(Pbm, Pbr) / (math.Pi * Dcp * data.Gasket.Width)
	}

	if bolt.SigmaB1 <= bolt.DSigmaM {
		bolt.VSigmaB1 = true
	}
	if bolt.SigmaB2 <= bolt.DSigmaR {
		bolt.VSigmaB2 = true
	}

	return bolt
}

// Расчет момента затяжки (иногда требуется не полных расчет, а только 2 значения. для определения подобных ситуаций используется флаг fullCalculate)
func (s *FlangeService) momentCalculate(
	Friction float64,
	data models.DataFlange,
	SigmaB1, DSigmaM, Pbm, Ab, Dcp float64,
	fullCalculate bool,
) *flange_model.CalcMoment {
	moment := &flange_model.CalcMoment{
		Friction: Friction,
	}

	if SigmaB1 > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter && data.Bolt.Diameter <= constants.MaxDiameter {
		moment.Mkp = s.graphic.CalculateMkp(data.Bolt.Diameter, SigmaB1)
		moment.UseGraphic = true
	} else {
		//? вроде как формула изменилась, но почему-то использовалась новая формула
		moment.Mkp = (Friction * Pbm * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
		// moment.Mkp = (0.3 * Pbm * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
	}

	if Friction == constants.DefaultFriction {
		moment.Mkp1 = 0.75 * moment.Mkp
	}

	if fullCalculate {
		Prek := 0.8 * Ab * data.Bolt.SigmaAt20
		moment.Qrek = Prek / (math.Pi * Dcp * data.Gasket.Width)
		moment.Mrek = (Friction * Prek * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
		// moment.Mrek = (0.3 * Prek * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000

		Pmax := DSigmaM * Ab
		moment.Qmax = Pmax / (math.Pi * Dcp * data.Gasket.Width)

		// Лукиных попросил вернуть условие 09.01.2024
		if data.TypeGasket == flange_model.GasketData_Soft && moment.Qmax > data.Gasket.PermissiblePres {
			Pmax = float64(data.Gasket.PermissiblePres) * (math.Pi * Dcp * data.Gasket.Width)
			moment.Qmax = data.Gasket.PermissiblePres
		}
		// if moment.Qmax > data.Gasket.PermissiblePres {
		// 	Pmax = float64(data.Gasket.PermissiblePres) * (math.Pi * Dcp * data.Gasket.Width)
		// 	moment.Qmax = data.Gasket.PermissiblePres
		// }

		moment.Mmax = (Friction * Pmax * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
		// moment.Mmax = (0.3 * Pmax * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000

		if moment.Mrek > moment.Mmax {
			moment.Mrek = moment.Mmax
		}
		if moment.Qrek > moment.Qmax {
			moment.Qrek = moment.Qmax
		}
	}

	return moment
}
