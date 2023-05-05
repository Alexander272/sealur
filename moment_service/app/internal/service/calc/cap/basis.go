package cap

import (
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/cap_model"
)

// Расчет основных величин
func (s *CapService) basisCalculate(data models.DataCap, req *calc_api.CapRequest) (*cap_model.Calculated_Basis, *cap_model.CalcAuxiliary) {
	deformation := s.deformationCalculate(data, req.Data.Pressure)
	forces, aux := s.forcesInBoltsCalculate(data, deformation, req)
	bolts := s.boltStrengthCalculate(data, req, forces.Pb, forces.Pbr, forces.A, deformation.Dcp, true)
	moment := &cap_model.CalcMoment{}

	ok := (bolts.VSigmaB1 && bolts.VSigmaB2 && data.TypeGasket != cap_model.GasketData_Soft) ||
		(bolts.VSigmaB1 && bolts.VSigmaB2 && bolts.Q <= float64(data.Gasket.PermissiblePres) && data.TypeGasket == cap_model.GasketData_Soft)
	if ok {
		moment = s.momentCalculate(req.Data.Friction, data, bolts.SigmaB1, bolts.DSigmaM, forces.Pb, forces.A, deformation.Dcp, true)
	}

	res := &cap_model.Calculated_Basis{
		Deformation:   deformation,
		ForcesInBolts: forces,
		BoltStrength:  bolts,
		Moment:        moment,
	}

	return res, aux
}

// Усилия, необходимые для смятия прокладки и обеспечения герметичности фланцевого соединения
func (s *CapService) deformationCalculate(data models.DataCap, pressure float64) *cap_model.CalcDeformation {
	deformation := &cap_model.CalcDeformation{}

	if data.TypeGasket == cap_model.GasketData_Oval {
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
func (s *CapService) forcesInBoltsCalculate(
	data models.DataCap,
	def *cap_model.CalcDeformation,
	req *calc_api.CapRequest,
) (*cap_model.CalcForcesInBolts, *cap_model.CalcAuxiliary) {
	forces := &cap_model.CalcForcesInBolts{}

	// формула 8
	// Суммарная площадь сечения болтов/шпилек
	forces.A = float64(data.Bolt.Count) * data.Bolt.Area

	// формула 9
	// Равнодействующая нагрузка от давления
	forces.Qd = 0.785 * math.Pow(def.Dcp, 2) * req.Data.Pressure

	// формула 10
	// Приведенная нагрузка, вызванная воздействием внешней силы и изгибающего момента
	forces.Qfm = float64(req.Data.AxialForce)

	var yp float64 = 0
	if data.TypeGasket == cap_model.GasketData_Soft {
		// Податливость прокладки
		yp = (data.Gasket.Thickness * data.Gasket.Compression) / (data.Gasket.Epsilon * math.Pi * def.Dcp * data.Gasket.Width)
	}

	// приложение К пояснение к формуле К.2
	Lb := data.Bolt.Length + s.typeBolt[req.Data.Type.String()]*data.Bolt.Diameter
	// формула К.2
	// Податливость болтов/шпилек
	yb := Lb / (data.Bolt.EpsilonAt20 * data.Bolt.Area * float64(data.Bolt.Count))

	flange := s.auxFlangeCalculate(req.FlangeData.Type, data.Flange, def.Dcp)
	cap := s.auxCapCalculate(req.CapData.Type, data.Cap, data.Flange, def.Dcp)

	aux := &cap_model.CalcAuxiliary{
		Flange: flange,
		Cap:    cap,
	}
	aux.Yp = yp
	aux.Lb = Lb
	aux.Yb = yb

	if data.TypeGasket == cap_model.GasketData_Oval || data.FlangeType == cap_model.FlangeData_free {
		// Для фланцев с овальными и восьмигранными прокладками и для свободных фланцев коэффициенты жесткости фланцевого соединения принимают равными 1.
		forces.Alpha = 1
	} else {
		// формула (Е.11)
		// Коэффициент жесткости
		forces.Alpha = 1 - (yp-(flange.Yf*flange.E+cap.Y*flange.B)*flange.B)/(yp+yb+(flange.Yf+cap.Y)*math.Pow(flange.B, 2))
	}

	minB := 0.4 * forces.A * data.Bolt.SigmaAt20
	forces.MinB = minB
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения обжатия прокладки и минимального начального натяжения болтов/шпилек
	forces.Pb2 = math.Max(def.Po, minB)
	// Расчетная нагрузка на болты/шпильки при затяжке, необходимая для обеспечения в рабочих условиях давления на
	// прокладку достаточного для герметизации фланцевого соединения
	forces.Pb1 = forces.Alpha*(forces.Qd+float64(req.Data.AxialForce)) + def.Rp

	divider := yp + yb*data.Bolt.EpsilonAt20/data.Bolt.Epsilon + (flange.Yf*data.Flange.EpsilonAt20/data.Flange.Epsilon)*math.Pow(flange.B, 2) +
		(cap.Y*data.Cap.EpsilonAt20/data.Cap.Epsilon)*math.Pow(flange.B, 2)

	if data.FlangeType == cap_model.FlangeData_free {
		divider += (flange.Yk * data.Flange.Ring.EpsilonAt20 / data.Flange.Ring.Epsilon) * math.Pow(flange.A, 2)
	}

	// формула (Е.8)
	gamma := 1 / divider
	aux.Gamma = gamma

	var temp1, temp2 float64
	if req.IsUseWasher {
		temp1 = (data.Flange.Alpha*data.Flange.H+data.Washer1.Alpha*data.Washer1.Thickness)*(data.Flange.T-20) +
			(data.Cap.Alpha*data.Cap.H+data.Washer2.Alpha*data.Washer2.Thickness)*(data.Cap.T-20)
	} else {
		temp1 = data.Flange.Alpha*data.Flange.H*(data.Flange.T-20) + data.Cap.Alpha*data.Cap.H*(data.Cap.T-20)
	}
	temp2 = data.Flange.H + data.Flange.H

	if data.FlangeType == cap_model.FlangeData_free {
		temp1 += data.Flange.Ring.Alpha * data.Flange.H * (data.Flange.Ring.T - 20)
		temp2 += data.Flange.Ring.Hk
	}
	if req.Data.IsEmbedded {
		temp1 += data.Embed.Alpha * data.Embed.Thickness * (req.Data.Temp - 20)
		temp2 += data.Embed.Thickness
	}

	//? должно быть два варианта формулы с шайбой и без нее
	// шайба будет задаваться так же как и болты + толщина шайбы

	//формула 11 (в старом 13)
	Qt := gamma * (temp1 - data.Bolt.Alpha*temp2*(data.Bolt.Temp-20))
	forces.Qt = Qt

	forces.Pb1 = math.Max(forces.Pb1, forces.Pb1-Qt)
	forces.Pb = math.Max(forces.Pb1, forces.Pb2)
	forces.Pbr = forces.Pb + (1-forces.Alpha)*(forces.Qd+float64(req.Data.AxialForce)) + Qt

	return forces, aux
}

// Проверка прочности болтов (шпилек) и прокладки
func (s *CapService) boltStrengthCalculate(
	data models.DataCap,
	req *calc_api.CapRequest,
	Pbm, Pbr, Ab, Dcp float64,
	isLoad bool,
) *cap_model.CalcBoltStrength {
	bolt := &cap_model.CalcBoltStrength{}

	bolt.SigmaB1 = Pbm / Ab
	bolt.SigmaB2 = Pbr / Ab

	Kyp := s.Kyp[req.Data.IsWork]
	Kyz := s.Kyz[req.Data.Condition.String()]
	Kyt := s.Kyt[isLoad]
	// формула Г.3
	bolt.DSigmaM = 1.2 * Kyp * Kyz * Kyt * data.Bolt.SigmaAt20
	// формула Г.4
	bolt.DSigmaR = Kyp * Kyz * Kyt * data.Bolt.Sigma

	if data.TypeGasket == cap_model.GasketData_Soft {
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
func (s *CapService) momentCalculate(
	Friction float64,
	data models.DataCap,
	SigmaB1, DSigmaM, Pbm, Ab, Dcp float64,
	fullCalculate bool,
) *cap_model.CalcMoment {
	moment := &cap_model.CalcMoment{
		Friction: Friction,
	}

	if SigmaB1 > constants.MaxSigmaB && data.Bolt.Diameter >= constants.MinDiameter && data.Bolt.Diameter <= constants.MaxDiameter {
		//TODO возвращать Friction и как-то определять на клиенте считается ли по формуле или по графику
		moment.Mkp = s.graphic.CalculateMkp(data.Bolt.Diameter, SigmaB1)
		moment.UseGraphic = true
	} else {
		//? вроде как формула изменилась, но почему-то использовалась новая формула
		moment.Mkp = (Friction * Pbm * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
		// moment.Mkp = (0.3 * Pbm * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
	}

	moment.Mkp1 = 0.75 * moment.Mkp

	if fullCalculate {
		Prek := 0.8 * Ab * data.Bolt.SigmaAt20
		moment.Qrek = Prek / (math.Pi * Dcp * data.Gasket.Width)
		moment.Mrek = (Friction * Prek * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000
		// moment.Mrek = (0.3 * Prek * data.Bolt.Diameter / float64(data.Bolt.Count)) / 1000

		Pmax := DSigmaM * Ab
		moment.Qmax = Pmax / (math.Pi * Dcp * data.Gasket.Width)

		if data.TypeGasket == cap_model.GasketData_Soft && moment.Qmax > data.Gasket.PermissiblePres {
			Pmax = float64(data.Gasket.PermissiblePres) * (math.Pi * Dcp * data.Gasket.Width)
			moment.Qmax = data.Gasket.PermissiblePres
		}

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
