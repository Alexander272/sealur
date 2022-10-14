package dev_cooling

import (
	"context"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/service/calc/dev_cooling/data"
	"github.com/Alexander272/sealur/moment_service/internal/service/flange"
	"github.com/Alexander272/sealur/moment_service/internal/service/gasket"
	"github.com/Alexander272/sealur/moment_service/internal/service/graphic"
	"github.com/Alexander272/sealur/moment_service/internal/service/materials"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

type CoolingService struct {
	graphic *graphic.GraphicService
	data    *data.DataService
	// formulas *formulas.FormulasService
	typeBolt map[string]float64
	mu       map[string]float64
}

func NewCoolingService(graphic *graphic.GraphicService, flange *flange.FlangeService, gasket *gasket.GasketService,
	materials *materials.MaterialsService) *CoolingService {
	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}

	mu := map[string]float64{
		"flat":   constants.Flat,
		"groove": constants.Groove,
	}

	data := data.NewDataService(flange, materials, gasket, graphic)
	// formulas := formulas.NewFormulasService()

	return &CoolingService{
		typeBolt: bolt,
		mu:       mu,
		graphic:  graphic,
		data:     data,
		// formulas: formulas,
	}
}

// расчет аппаратов воздушного охлаждения по ГОСТ 25822-83
func (s *CoolingService) CalculateDevCooling(ctx context.Context, data *calc_api.DevCoolingRequest) (*calc_api.DevCoolingResponse, error) {
	d, err := s.data.GetData(ctx, data)
	if err != nil {
		return nil, err
	}

	result := calc_api.DevCoolingResponse{
		Data:      s.data.FormatInitData(data),
		Cap:       d.Cap,
		TubeSheet: d.TubeSheet,
		Tube:      d.Tube,
		Bolts:     d.Bolt,
		Gasket:    d.Gasket,
		Calc:      &dev_cooling_model.Calculated{},
		Formulas:  &dev_cooling_model.Formulas{},
	}

	Auxiliary := &dev_cooling_model.CalcAuxiliary{}
	Bolt := &dev_cooling_model.CalcBolt{}
	TubeSheet := &dev_cooling_model.CalcTubeSheet{}
	Cap := &dev_cooling_model.CalcCap{}

	// расчетная ширина плоской прокладки
	Auxiliary.EstimatedGasketWidth = math.Min(d.Gasket.Width, 3.87*math.Sqrt(d.Gasket.Width))
	// расчетный размер решетки в поперечном направлении
	Auxiliary.Bp = d.Gasket.SizeTrans - Auxiliary.EstimatedGasketWidth

	// Условия применения формул
	result.Calc.Condition1 = &dev_cooling_model.Condition{
		X: (d.TubeSheet.ZoneThick - d.TubeSheet.Corrosion) / Auxiliary.Bp,
		Y: constants.Condition,
	}
	result.Calc.Condition2 = &dev_cooling_model.Condition{
		X: (d.Cap.BottomThick - d.Cap.Corrosion) / Auxiliary.Bp,
		Y: constants.Condition,
	}

	if result.Calc.Condition1.X > result.Calc.Condition1.Y || result.Calc.Condition2.X > result.Calc.Condition2.Y {
		return &result, nil
	}
	result.Calc.IsConditionsMet = true

	cond1 := math.Min(d.Cap.SigmaAt20/d.Cap.Sigma, d.TubeSheet.SigmaAt20/d.TubeSheet.Sigma)
	cond2 := math.Min(d.Tube.SigmaAt20/d.Tube.Sigma, d.Bolt.SigmaAt20/d.Bolt.Sigma)

	// Пробное давление
	result.Calc.Pressure = 1.25 * data.Pressure * math.Min(cond1, cond2)

	// расчетная ширина перфорированной зоны решетки
	Auxiliary.EstimatedZoneWidth = math.Min(float64(d.TubeSheet.Count)*d.TubeSheet.StepTrans, Auxiliary.Bp)
	// относительная ширина беструбной зоны решетки
	Auxiliary.RelativeWidth = (Auxiliary.Bp - Auxiliary.EstimatedZoneWidth) / Auxiliary.EstimatedZoneWidth

	// Вспомогательные коэффициенты
	Auxiliary.Upsilon = (math.Pi * (d.Tube.Diameter - d.Tube.Thickness) * (d.Tube.Thickness - d.Tube.Corrosion)) /
		(d.TubeSheet.StepLong * d.TubeSheet.StepTrans)
	Auxiliary.Eta = 1 - (math.Pi/4)*(math.Pow(d.Tube.Diameter-2*d.Tube.Thickness, 2)/(d.TubeSheet.StepLong*d.TubeSheet.StepTrans))

	// эффективный диаметр отверстия решетки или задней стенке
	if data.Method == calc_api.DevCoolingRequest_AllThickness {
		Auxiliary.D = d.TubeSheet.Diameter - 2*d.Tube.Thickness
	} else if data.Method == calc_api.DevCoolingRequest_PartThickness {
		Auxiliary.D = d.TubeSheet.Diameter - d.Tube.Thickness
	} else {
		Auxiliary.D = d.TubeSheet.Diameter
	}

	// коэффициент ослабления решетки и задней стенки
	Auxiliary.Phi = 1 - Auxiliary.D/d.TubeSheet.StepLong
	// допускаемая нагрузка из условия прочности труб
	Auxiliary.LoadTube = Auxiliary.Upsilon * (1 - ((d.Tube.Diameter-d.Tube.Thickness)/(2*(d.Tube.Thickness-d.Tube.Corrosion)))*
		(data.Pressure/d.Tube.Sigma)) * d.Tube.Sigma

	Auxiliary.Mu = s.mu[data.TypeMounting.String()]

	// допускаемое напряжение из условия прочности крепления трубы в решетке
	if data.Mounting == calc_api.DevCoolingRequest_flaring {
		Auxiliary.Load = Auxiliary.Upsilon * Auxiliary.Mu * ((2 * d.Tube.Depth) / (d.Tube.Diameter * d.Tube.Thickness)) * d.Tube.Sigma
	} else if data.Mounting == calc_api.DevCoolingRequest_welding {
		Auxiliary.Load = 0.7 * Auxiliary.Upsilon * (d.Tube.Size / d.Tube.Thickness) * math.Min(d.Tube.Sigma, d.TubeSheet.Sigma)
	} else {
		Auxiliary.Load = 0.7*Auxiliary.Upsilon*(d.Tube.Size/d.Tube.Thickness)*math.Min(d.Tube.Sigma, d.TubeSheet.Sigma) +
			0.6*(Auxiliary.Upsilon*Auxiliary.Mu*((2*d.Tube.Depth)/(d.Tube.Diameter*d.Tube.Thickness))*d.Tube.Sigma)
	}

	// коэффициент уменьшения допускаемых напряжений при продольном изгибе
	//! тут в Epsilon лежит значение EpsilonAt20
	Auxiliary.PhiT = 1 / math.Sqrt(1+math.Pow(1.8*(d.Tube.Sigma/d.Tube.Epsilon)*
		math.Pow(d.Tube.ReducedLength/(d.Tube.Diameter-d.Tube.Thickness), 2), 2))

	// l1, l2 - Плечи изгибающих моментов
	Auxiliary.Arm1 = 0.5 * (d.Bolt.Distance - Auxiliary.Bp)
	Auxiliary.Arm2 = 0.5 * (d.Bolt.Distance - d.Gasket.SizeTrans)

	//? Если 5 чертеж то wallThick (s5) будет равна 0
	if data.CameraDiagram == calc_api.DevCoolingRequest_schema5 {
		d.Cap.WallThick = d.Cap.BottomThick
		d.Cap.Radius = d.Cap.InnerSize / 2
	}

	// Phi для Угловые податливости крышки
	var Phi1, Phi2, Phi3, Phi4, Phi5, Phi6 float64
	if data.CameraDiagram == calc_api.DevCoolingRequest_schema1 || data.CameraDiagram == calc_api.DevCoolingRequest_schema4 {
		Phi1 = 1
		Phi2 = 8 * math.Pow(d.Cap.Depth/d.Cap.InnerSize, 3)
		Phi4 = 1
		Phi5 = 2 * (d.Cap.Depth / d.Cap.InnerSize)
	} else {
		Phi1 = 1 + 0.85*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 2) - 12.55*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 3) +
			13.7*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 2)*(d.Cap.Depth/d.Cap.InnerSize)
		Phi2 = 8*math.Pow(d.Cap.Depth/d.Cap.InnerSize, 3) - 12*(d.Cap.Depth/d.Cap.InnerSize)*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 2) +
			4*math.Pow(d.Cap.Radius/d.Cap.InnerSize, 3)
		Phi4 = 1 - 1.14*(d.Cap.Radius/d.Cap.InnerSize)
		Phi5 = 2*(d.Cap.Depth/d.Cap.InnerSize) - 2*(d.Cap.Radius/d.Cap.InnerSize)
	}
	Phi3 = 12*math.Pow(d.Cap.Depth/d.Cap.InnerSize, 2)*(d.Cap.FlangeThick/d.Cap.InnerSize) - 4*math.Pow(d.Cap.FlangeThick/d.Cap.InnerSize, 3)
	Phi6 = 2 * (d.Cap.FlangeThick / d.Cap.InnerSize)

	Cap.Psi = (math.Pow(Auxiliary.Bp/d.Cap.InnerSize, 2)-1)*(d.Cap.L/(d.Cap.L+d.Cap.InnerSize)) - 4*math.Pow(d.Cap.Depth/d.Cap.InnerSize, 2)
	Lb := d.Cap.L / d.Cap.InnerSize

	var eta float64
	if data.CameraDiagram == calc_api.DevCoolingRequest_schema4 {
		eta = 0.275 * (math.Pow(d.Cap.OuterSize-d.Cap.InnerSize, 4) / d.Cap.L * math.Pow(d.Cap.BottomThick, 4)) *
			math.Pow((2*d.Cap.FlangeThick)/(d.Cap.OuterSize-d.Cap.InnerSize)-0.63, 2)
	} else {
		eta = 4.4 * (d.Cap.OuterSize / d.Cap.L) * (((2*d.Cap.FlangeThick)/(d.Cap.OuterSize-d.Cap.InnerSize)-0.63)*math.Pow(
			d.Cap.FlangeThick/d.Cap.BottomThick, 4) + ((d.Cap.Depth+d.Cap.WallThick-d.Cap.FlangeThick)/
			d.Cap.BottomThick-0.63)*math.Pow(d.Cap.WallThick/d.Cap.BottomThick, 4))
	}

	var x1, x2, m float64 = 0, 0, 1
	for {
		alphaM := m * math.Pi * d.Cap.InnerSize / (2 * d.Cap.L)

		part1 := 1 / math.Pow(m, 3)
		part2 := alphaM - (alphaM*math.Tanh(alphaM)+1)*math.Tanh(alphaM)
		part3 := eta * (alphaM - (alphaM*math.Tanh(alphaM)-1)*math.Tanh(alphaM))
		part4 := 2 + m*eta*(alphaM-(alphaM*math.Tanh(alphaM)-1)*math.Tanh(alphaM))
		part5 := math.Sin(math.Pi * m / 2)

		x2 = part1 * (part2 * ((part3 / part4) - (1 / m))) * part5

		x1 += x2
		m++

		if m > 2 && (x2/x1)*1000 >= 0.0001 {
			break
		}
	}
	Lambda1 := 0.0206 * math.Pow(Lb, 3) * x1

	x1, x2, m = 0, 0, 1
	for {
		alphaM := m * math.Pi * d.Cap.InnerSize / (2 * d.Cap.L)

		part1 := alphaM - (alphaM*math.Tanh(alphaM)-1)*math.Tanh(alphaM)
		part2 := math.Sin(math.Pi * m / 2)
		part3 := math.Pow(m, 2) * (2 + m*eta*(alphaM-(alphaM*math.Tanh(alphaM)-1)*math.Tanh(alphaM)))

		x2 = (part1 * part2) / part3

		x1 += x2
		m++

		if m > 2 && (x2/x1)*1000 >= 0.0001 {
			break
		}
	}
	Lambda2 := 0.406 * Lb * x1

	tmp1 := math.Pow(d.Cap.InnerSize, 3) / (d.Cap.Epsilon * math.Pow(d.Cap.BottomThick, 3))
	var tmp2, tmp3, tmpPhi float64
	if data.CameraDiagram == calc_api.DevCoolingRequest_schema4 {
		tmpPhi = Phi4
		tmp2 = Phi1 * Lambda1
	} else {
		tmp2 = (Phi1 + (Phi2-Phi3)*math.Pow(d.Cap.BottomThick/d.Cap.WallThick, 3)) * Lambda1
		tmpPhi = Phi4 + (Phi5-Phi6)*math.Pow(d.Cap.BottomThick/d.Cap.WallThick, 3)
	}
	tmp3 = (1.0 / 8) * tmpPhi * Cap.Psi * Lambda2
	// Угловые податливости крышки
	Bolt.CapUpsilonP = 10.9 * tmp1 * (tmp2 + tmp3)

	Bolt.Lp = d.Gasket.SizeLong - Auxiliary.EstimatedGasketWidth
	// Угловые податливости крышки
	Bolt.CapUpsilonM = 10.9 * (d.Cap.InnerSize / (2 * d.Cap.Epsilon * math.Pow(d.Cap.BottomThick, 3) * (Bolt.Lp + Auxiliary.Bp))) * tmpPhi * Lambda2

	TubeSheet.Psi = Auxiliary.RelativeWidth * (Auxiliary.RelativeWidth + 2)
	TubeSheet.Omega = 1.6 * (Auxiliary.EstimatedZoneWidth / d.TubeSheet.ZoneThick) * math.Pow(Auxiliary.Upsilon*d.TubeSheet.ZoneThick/d.Tube.Length, 1.0/4)

	Y1 := math.Cosh(TubeSheet.Omega) * math.Cos(TubeSheet.Omega)
	Y2 := 0.5 * (math.Cosh(TubeSheet.Omega)*math.Sin(TubeSheet.Omega) + math.Sinh(TubeSheet.Omega)*math.Cos(TubeSheet.Omega))
	Y3 := 0.5 * math.Sinh(TubeSheet.Omega) * math.Sin(TubeSheet.Omega)
	Y4 := 0.25 * (math.Cosh(TubeSheet.Omega)*math.Sin(TubeSheet.Omega) - math.Sinh(TubeSheet.Omega)*math.Cos(TubeSheet.Omega))

	// $alfa1 = ($Y2 - $Y2 * $Y1 - 4 * $Y4 * $Y3) / ($omega * ($Y2 * $Y4 - $Y3 * $Y3));
	Alpha1 := (Y2 - Y2*Y1 - 4*Y4*Y3) / (TubeSheet.Omega * (Y2*Y4 - math.Pow(Y3, 2)))
	// $alfa2 = ($Y1 * $Y3 + $Y3 - $Y2 * $Y2) / ($omega * $omega * ($Y2 * $Y4 - $Y3 * $Y3));
	Alpha2 := (Y1*Y3 + Y3 - math.Pow(Y2, 2)) / (math.Pow(TubeSheet.Omega, 2) * (Y2*Y4 - math.Pow(Y3, 2)))

	tmp1 = 0.23 * math.Pow(Auxiliary.EstimatedZoneWidth, 3) / (d.TubeSheet.Epsilon * math.Pow(d.TubeSheet.ZoneThick, 3))
	tmp2 = Auxiliary.RelativeWidth * (2*TubeSheet.Psi - Auxiliary.RelativeWidth) * math.Pow(d.TubeSheet.ZoneThick/d.TubeSheet.OutZoneThick, 3)
	tmp3 = 1.7 * (TubeSheet.Psi*Alpha1 + 4*Alpha2)
	// Угловые податливости решетки
	Bolt.SheetUpsilonP = tmp1 * (tmp2 + tmp3)

	tmp1 = Auxiliary.EstimatedZoneWidth / (2 * d.TubeSheet.Epsilon * math.Pow(d.TubeSheet.ZoneThick, 3) * (Bolt.Lp + Auxiliary.Bp))
	tmp2 = 2*Auxiliary.RelativeWidth*math.Pow(d.TubeSheet.ZoneThick/d.TubeSheet.OutZoneThick, 3) + 1.1*Alpha1
	// Угловые податливости решетки
	Bolt.SheetUpsilonM = 2.7 * tmp1 * tmp2

	Lb = d.Bolt.Lenght + s.typeBolt[data.TypeBolt.String()]*d.Bolt.Diameter

	// Линейная податливость болта (шпильки)
	Bolt.UpsilonB = Lb / (d.Bolt.Epsilon * d.Bolt.Area * float64(d.Bolt.Count))
	// Линейная податливость прокладки
	Bolt.UpsilonP = d.Gasket.Thickness / (2 * d.Bolt.Epsilon * (Bolt.Lp + Auxiliary.Bp) * d.Gasket.Width)

	tmp1 = Bolt.UpsilonB + (Bolt.CapUpsilonM+Bolt.SheetUpsilonM)*math.Pow(Auxiliary.Arm1, 2)
	tmp2 = ((Bolt.CapUpsilonP + Bolt.SheetUpsilonP) / (Bolt.Lp * Auxiliary.Bp)) * Auxiliary.Arm1
	tmp3 = Bolt.UpsilonB + Bolt.UpsilonP + (Bolt.CapUpsilonM+Bolt.SheetUpsilonM)*math.Pow(Auxiliary.Arm1, 2)
	// Коэффициент податливости фланцевого соединения крышки и решетки
	Bolt.Eta = (tmp1 + tmp2) / tmp3

	// Fв - Расчетное усилие в болтах (шпильках) в условиях эксплуатации
	Bolt.WorkEffort = data.Pressure * (Bolt.Lp*Auxiliary.Bp + 2*Auxiliary.EstimatedGasketWidth*d.Gasket.M*(Bolt.Lp+Auxiliary.Bp))

	//TODO в оригинале почему-то тут не WorkEffort, а полщадь и количество болтов
	tmp1 = (result.Calc.Pressure / data.Pressure) * Bolt.WorkEffort
	tmp2 = result.Calc.Pressure * (Bolt.Eta*Bolt.Lp*Auxiliary.Bp + 2*Auxiliary.EstimatedGasketWidth*d.Gasket.M*(Bolt.Lp+Auxiliary.Bp))
	// F0 - Расчетное усилие в болтах (шпильках) в условиях испытаний или монтажа
	Bolt.Effort = math.Max(tmp1, tmp2)

	Ab := d.Bolt.Area * float64(d.Bolt.Count)
	// Условия прочности болтов/шпилек - в условиях испытания или монтажа
	Bolt.TestCond = &dev_cooling_model.Condition{
		X: Bolt.Effort / Ab,
		Y: d.Bolt.SigmaAt20,
	}
	// Условия прочности болтов/шпилек - в условиях эксплуатации
	Bolt.WorkCond = &dev_cooling_model.Condition{
		X: Bolt.WorkEffort / Ab,
		Y: d.Bolt.Sigma,
	}

	result.Calc.Auxiliary = Auxiliary
	result.Calc.Bolt = Bolt
	result.Calc.TubeSheet = TubeSheet
	result.Calc.Cap = Cap

	return &result, nil
}
