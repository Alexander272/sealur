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
	Auxiliary.EstimatedZoneWidth = math.Min(d.Gasket.Width, 3.87*math.Sqrt(d.Gasket.Width))
	// расчетный размер решетки в поперечном направлении
	Auxiliary.Bp = d.Gasket.SizeTrans - Auxiliary.EstimatedZoneWidth

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
	// TODO похоже тут EpsilonAt20 нужен, а не просто Epsilon
	Auxiliary.PhiT = 1 / math.Sqrt(1+math.Pow(1.8*(d.Tube.Sigma/d.Tube.Epsilon)*
		math.Pow(d.Tube.ReducedLength/(d.Tube.Diameter-d.Tube.Thickness), 2), 2))

	// l1, l2 - Плечи изгибающих моментов
	Auxiliary.Arm1 = 0.5 * (d.Bolt.Distance - Auxiliary.Bp)
	Auxiliary.Arm2 = 0.5 * (d.Bolt.Distance - d.Gasket.SizeTrans)

	result.Calc.Auxiliary = Auxiliary
	result.Calc.Bolt = Bolt
	result.Calc.TubeSheet = TubeSheet
	result.Calc.Cap = Cap

	return &result, nil
}
