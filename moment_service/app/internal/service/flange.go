package service

import (
	"context"
	"fmt"
	"math"

	"github.com/Alexander272/sealur/moment_service/internal/constants"
	"github.com/Alexander272/sealur/moment_service/internal/models"
	"github.com/Alexander272/sealur/moment_service/internal/repository"
	moment_proto "github.com/Alexander272/sealur/moment_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/moment_service/pkg/logger"
)

type FlangeService struct {
	repo          repository.Flange
	materials     *MaterialsService
	gasket        *GasketService
	typeFlangesTF map[string]float64
	typeFlangesTD map[string]float64
	typeBolt      map[string]float64
}

func NewFlangeService(repo repository.Flange, materials *MaterialsService, gasket *GasketService) *FlangeService {
	flangesTF := map[string]float64{
		"isolated":    constants.IsolatedFlatTf,
		"nonIsolated": constants.NonIsolatedFlatTf,
	}

	flangesTD := map[string]float64{
		"isolated":         constants.IsolatedFlatTb,
		"nonIsolated":      constants.NonIsolatedFlatTb,
		"isolated-free":    constants.IsolatedFreeTb,
		"nonIsolated-free": constants.NonIsolatedFlatTb,
	}

	bolt := map[string]float64{
		"bolt": constants.BoltD,
		"pin":  constants.PinD,
	}

	return &FlangeService{
		repo:          repo,
		materials:     materials,
		gasket:        gasket,
		typeFlangesTF: flangesTF,
		typeFlangesTD: flangesTD,
		typeBolt:      bolt,
	}
}

// TODO в зависимости от госта можно будет вызывать отдельные функции
func (s *FlangeService) Calculation(ctx context.Context, data *moment_proto.FlangeRequest) (*moment_proto.FlangeResponse, error) {
	dataFlangeFirst, err := s.getDataFlange(ctx, data.FlangesData[0], data.Flanges, data.Temp)
	if err != nil {
		return nil, err
	}
	var dataFlangeSecond models.InitialDataFlange
	if len(data.FlangesData) > 1 {
		dataFlangeSecond, err = s.getDataFlange(ctx, data.FlangesData[1], data.Flanges, data.Temp)
		if err != nil {
			return nil, err
		}
	} else {
		dataFlangeSecond = dataFlangeFirst
	}

	//TODO добавить зависимость от типа фланца
	Tb := s.typeFlangesTD[data.Flanges] * float64(data.Temp)
	// if ($TipF1 == 2) {

	//TODO
	boltMat, err := s.materials.GetMatFotCalculate(ctx, data.Bolt.MarkId, Tb)
	if err != nil {
		return nil, err
	}
	logger.Debug(boltMat)

	gasket, err := s.gasket.Get(ctx, models.GetGasket{TypeGasket: "", Env: "", Thickness: 3.2})
	if err != nil {
		return nil, err
	}
	var Lb0 float64
	// TODO здесь какое-то условие надо
	Lb0 = float64(gasket.Thickness)
	Lb0 += dataFlangeFirst.H + dataFlangeSecond.H

	logger.Debug(gasket)

	bp := float64(data.Gasket.DOut-data.Gasket.DIn) / 2

	var b0, Dcp float64

	if gasket.IsOval {
		b0 = bp / 4
		Dcp = float64(data.Gasket.DOut) - bp/2
	} else {
		if bp <= constants.Bp {
			b0 = bp
		} else {
			b0 = constants.B0 * math.Sqrt(bp)
		}
		Dcp = float64(data.Gasket.DOut) - b0
	}

	// yp = 0.0
	// if ($TipP == 0) {
	// 	$yp = $hp * $Kp / ($Ep * pi() * $Dcp * $bp)
	// }
	var yp float64 = 0
	//? иногда она не считается
	yp = float64(gasket.Thickness*gasket.Compression) / (float64(gasket.Epsilon) * math.Pi * Dcp * bp)

	Lb := Lb0 + s.typeBolt[data.Type]*float64(dataFlangeFirst.Diameter)

	yb := Lb / (dataFlangeFirst.EpsilonAt20 * float64(dataFlangeFirst.Area) * float64(dataFlangeFirst.Count))
	Ab := float64(dataFlangeFirst.Count) * float64(dataFlangeFirst.Area)

	logger.Debug(yp, yb, Ab, b0)

	return &moment_proto.FlangeResponse{}, nil
}

// Функция для получения данных необходимых для расчетов
func (s *FlangeService) getDataFlange(
	ctx context.Context,
	flange *moment_proto.FlangeData,
	typeFlange string,
	temp float32,
) (models.InitialDataFlange, error) {

	size, err := s.repo.GetSize(ctx, float64(flange.Dy), float64(flange.Py))
	if err != nil {
		return models.InitialDataFlange{}, fmt.Errorf("failed to get size. error: %w", err)
	}

	var dataFlange models.InitialDataFlange

	//TODO добавить зависимость от типа фланца (можно добавить мапу из которой в зависимости от фланца будет выбираться нужная константа)
	dataFlange.Tf = s.typeFlangesTF[typeFlange] * float64(temp)
	// TODO тут еще Tk1 считается если фланцы свободные

	mat, err := s.materials.GetMatFotCalculate(ctx, flange.MarkId, dataFlange.Tf)
	if err != nil {
		return models.InitialDataFlange{}, err
	}

	dataFlange = models.InitialDataFlange{
		DOut:        size.D1,
		D:           size.D,
		H:           size.B,
		S0:          size.S0,
		S1:          size.S1,
		L:           size.Lenght,
		D6:          size.D2,
		Diameter:    size.Diameter,
		Count:       size.Count,
		Area:        size.Area,
		AlphaF:      mat.AlphaF,
		EpsilonAt20: mat.EpsilonAt20,
		Epsilon:     mat.Epsilon,
		SigmaAt20:   mat.SigmaAt20,
		Sigma:       mat.Sigma,
	}

	dataFlange.SigmaM = constants.SigmaM * mat.Sigma
	dataFlange.SigmaMAt20 = constants.SigmaM * mat.SigmaAt20
	dataFlange.SigmaR = constants.SigmaR * mat.Sigma
	dataFlange.SigmaRAt20 = constants.SigmaR * mat.SigmaAt20

	return dataFlange, nil
}

// расчеты
func (s *FlangeService) getCalculatedDate(
	ctx context.Context,
	flange *moment_proto.FlangeData,
	data models.InitialDataFlange,
	Dcp float64,
) (models.CalculatedData, error) {
	var calculated models.CalculatedData
	if flange.Type != "free" {
		calculated.B = 0.5 * (data.D6 - Dcp)
	} else {
		calculated.Ds = 0.5 * (data.D + data.Dk + 2*data.H)
		calculated.A = 0.5 * (data.D6 - data.Ds)
		calculated.B = 0.5 * (data.Ds - Dcp)
	}

	if flange.Type != "welded" {
		calculated.Se = data.S0
	} else {
		calculated.X = data.L / (math.Sqrt(data.D * data.S0))
		calculated.Betta = data.S1 / data.S0
		calculated.Zak = 1 + (calculated.Betta-1)*calculated.X/(calculated.X+(1+calculated.Betta)/4)
		calculated.Se = calculated.Zak * data.S0
	}

	calculated.E = 0.5 * (Dcp - data.DOut - calculated.Se)
	calculated.L0 = math.Sqrt(data.DOut * data.S0)
	calculated.K = data.D / data.DOut

	dividend := math.Pow(calculated.K, 2)*(1+8.55*(math.Log(calculated.K)/math.Log(10))) - 1
	divider := (1.05 + 1.945*math.Pow(calculated.K, 2)) * (calculated.K - 1)
	calculated.BettaT = dividend / divider

	divider = 1.36 * (math.Pow(calculated.K, 2) - 1) * (calculated.K - 1)
	calculated.BettaU = dividend / divider

	dividend = 1 / (calculated.K - 1)
	divider = 0.69 + 5.72*((math.Pow(calculated.K, 2)*(math.Log(calculated.K)/math.Log(10)))/(math.Pow(calculated.K, 2)-1))
	calculated.BettaY = dividend * divider

	dividend = math.Pow(calculated.K, 2) + 1
	divider = math.Pow(calculated.K, 2) - 1
	calculated.BettaZ = dividend / divider

	if flange.Type == "welded" && data.S0 != data.S1 {
		betta := data.S1 / data.S0
		x := data.L / calculated.L0

		calculated.BettaF = s.calculateBettaF(betta, x)
		calculated.BettaV = s.calculateBettaV(betta, x)
		calculated.F = s.calculateF(betta, x)
	} else {
		calculated.BettaF = 0.908920
		calculated.BettaV = 0.550103
		calculated.F = 1.0
	}

	calculated.Lymda = (calculated.BettaF+data.H+calculated.L0)/calculated.BettaT*data.L + calculated.BettaV*math.Pow(data.H, 3)/(calculated.BettaU*calculated.L0*math.Pow(data.S0, 2))
	calculated.Yf = (0.91 * calculated.BettaV) / (data.EpsilonAt20 * calculated.Lymda * math.Pow(data.S0, 2) * calculated.L0)

	if flange.Type == "free" {
		// $psik1 = 1.28 * (log($Dnk1 / $Dk1) / log(10));
		// $yk1 = 1.0 / ($Ek201 * $hk1 * $hk1 * $hk1 * $psik1);
		// calculated.Psik = 1.28 * (math.Log(data.))
	}

	if flange.Type != "free" {
		// $yfn1 = (pow((pi() / 4), 3)) * ($Db1 / ($E201 * $Dn1 * pow($h1, 3)));
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.D6 / (data.EpsilonAt20 * data.D * math.Pow(data.H, 3)))
	} else {
		calculated.Yfn = math.Pow(math.Pi/4, 3) * (data.Ds / (data.EpsilonAt20 * data.D * math.Pow(data.H, 3)))
		// calculated.Yfc = math.Pow(math.Pi/4, 3) * (data.D6 / ())
		// $yfn1 = (pow((pi() / 4), 3)) * ($Ds1 / ($E201 * $Dn1 * pow($h1, 3)));
		// $yfc1 = (pow((pi() / 4), 3)) * ($Db1 / ($Ek201 * $Dnk1 * pow($hk1, 3)));
	}

	return calculated, nil
}

// функция тупо скопирована из оригинала
// расчет аппроксимированной функции (функция разная и зависит от значения х. исходные значения Рисунок К.2 ГОСТ 34233.4-2017)
func (s *FlangeService) calculateBettaF(betta, x float64) float64 {
	var f, f1, f2 float64

	switch {
	case x >= 0 && x < 0.1:
		f1 = 0.908920
		f2 = (-0.000000685709774295162)*math.Pow(betta, 6) + 0.00000179042442916*math.Pow(betta, 5) + 0.00121342871946961*math.Pow(betta, 4) - 0.0156079520766816*math.Pow(betta, 3) + 0.0713852548204228*math.Pow(betta, 2) - 0.132033833830155*betta + 0.983961997348035
		f = ((x-0.0)/(0.10-0.0))*(f2-f1) + f1
	case x >= 0.1 && x < 0.2:
		f1 = -0.000000685709774295162*math.Pow(betta, 6) + 0.00000179042442916*math.Pow(betta, 5) + 0.00121342871946961*math.Pow(betta, 4) - 0.0156079520766816*math.Pow(betta, 3) + 0.0713852548204228*math.Pow(betta, 2) - 0.132033833830155*betta + 0.983961997348035
		f2 = -0.000186958718629171*math.Pow(betta, 6) + 0.00270260935854338*math.Pow(betta, 5) - 0.0132402054724494*math.Pow(betta, 4) + 0.0177314503113593*math.Pow(betta, 3) + 0.0496071668843001*math.Pow(betta, 2) - 0.163661061259418*betta + 1.01596699859481
		f = ((x-0.10)/(0.20-0.10))*(f2-f1) + f1
	case x >= 0.2 && x < 0.25:
		f2 = 0.00173545909976035*math.Pow(betta, 4) - 0.0229593274470913*math.Pow(betta, 3) + 0.109804849073091*math.Pow(betta, 2) - 0.221341587836258*betta + 1.04140676493752
		f1 = -0.000186958718629171*math.Pow(betta, 6) + 0.00270260935854338*math.Pow(betta, 5) - 0.0132402054724494*math.Pow(betta, 4) + 0.0177314503113593*math.Pow(betta, 3) + 0.0496071668843001*math.Pow(betta, 2) - 0.163661061259418*betta + 1.01596699859481
		f = ((x-0.20)/(0.25-0.20))*(f2-f1) + f1
	case x >= 0.25 && x < 0.3:
		f1 = 0.00173545909976035*math.Pow(betta, 4) - 0.0229593274470913*math.Pow(betta, 3) + 0.109804849073091*math.Pow(betta, 2) - 0.221341587836258*betta + 1.04140676493752
		f2 = 0.00193732076800754*math.Pow(betta, 4) - 0.026106230052389*math.Pow(betta, 3) + 0.128138534713705*math.Pow(betta, 2) - 0.268021034819988*betta + 1.07275733229217
		f = ((x-0.25)/(0.30-0.25))*(f2-f1) + f1
	case x >= 0.3 && x < 0.35:
		f2 = 0.00297469010544045*math.Pow(betta, 4) - 0.0392245779011808*math.Pow(betta, 3) + 0.187025725151586*math.Pow(betta, 2) - 0.379373192875038*betta + 1.13713660128899
		f1 = 0.00193732076800754*math.Pow(betta, 4) - 0.026106230052389*math.Pow(betta, 3) + 0.128138534713705*math.Pow(betta, 2) - 0.268021034819988*betta + 1.07275733229217
		f = ((x-0.30)/(0.35-0.30))*(f2-f1) + f1
	case x >= 0.35 && x < 0.4:
		f1 = 0.00297469010544045*math.Pow(betta, 4) - 0.0392245779011808*math.Pow(betta, 3) + 0.187025725151586*math.Pow(betta, 2) - 0.379373192875038*betta + 1.13713660128899
		f2 = 0.00309341530338601*math.Pow(betta, 4) - 0.0410675648368301*math.Pow(betta, 3) + 0.198159262817235*math.Pow(betta, 2) - 0.411556783898828*betta + 1.15989660189039
		f = ((x-0.35)/(0.40-0.35))*(f2-f1) + f1
	case x >= 0.4 && x < 0.45:
		f2 = 0.00317467263939489*math.Pow(betta, 4) - 0.0425118835549549*math.Pow(betta, 3) + 0.208444710817431*math.Pow(betta, 2) - 0.445219873114703*betta + 1.18481173816372
		f1 = 0.00309341530338601*math.Pow(betta, 4) - 0.0410675648368301*math.Pow(betta, 3) + 0.198159262817235*math.Pow(betta, 2) - 0.411556783898828*betta + 1.15989660189039
		f = ((x-0.40)/(0.45-0.40))*(f2-f1) + f1
	case x >= 0.45 && x < 0.5:
		f1 = 0.00317467263939489*math.Pow(betta, 4) - 0.0425118835549549*math.Pow(betta, 3) + 0.208444710817431*math.Pow(betta, 2) - 0.445219873114703*betta + 1.18481173816372
		f2 = 0.00306882111532593*math.Pow(betta, 4) - 0.042363393407285*math.Pow(betta, 3) + 0.214087116035717*math.Pow(betta, 2) - 0.472733357010185*betta + 1.20668840864142
		f = ((x-0.45)/(0.50-0.45))*(f2-f1) + f1
	case x >= 0.5 && x < 0.6:
		f2 = 0.00431853154006122*math.Pow(betta, 4) - 0.0569156296982365*math.Pow(betta, 3) + 0.27529595841789*math.Pow(betta, 2) - 0.593152205919012*betta + 1.27878144516105
		f1 = 0.00306882111532593*math.Pow(betta, 4) - 0.042363393407285*math.Pow(betta, 3) + 0.214087116035717*math.Pow(betta, 2) - 0.47273335701018*betta + 1.20668840864142
		f = ((x-0.50)/(0.60-0.50))*(f2-f1) + f1
	case x >= 0.6 && x < 0.7:
		f1 = 0.00431853154006122*math.Pow(betta, 4) - 0.0569156296982365*math.Pow(betta, 3) + 0.27529595841789*math.Pow(betta, 2) - 0.593152205919012*betta + 1.27878144516105
		f2 = 0.00431178957617191*math.Pow(betta, 4) - 0.0578175725367722*math.Pow(betta, 3) + 0.285866341996179*math.Pow(betta, 2) - 0.636970282296766*betta + 1.31307260935091
		f = ((x-0.60)/(0.70-0.60))*(f2-f1) + f1
	case x >= 0.7 && x < 0.8:
		f2 = 0.00393976776368533*math.Pow(betta, 4) - 0.0534852765914117*math.Pow(betta, 3) + 0.271834547060149*math.Pow(betta, 2) - 0.638055491315338*betta + 1.32441823957794
		f1 = 0.00431178957617191*math.Pow(betta, 4) - 0.0578175725367722*math.Pow(betta, 3) + 0.285866341996179*math.Pow(betta, 2) - 0.636970282296766*betta + 1.31307260935091
		f = ((x-0.70)/(0.80-0.70))*(f2-f1) + f1
	case x >= 0.8 && x < 0.9:
		f1 = 0.00393976776368533*math.Pow(betta, 4) - 0.0534852765914117*math.Pow(betta, 3) + 0.271834547060149*math.Pow(betta, 2) - 0.638055491315338*betta + 1.32441823957794
		f2 = 0.00363323914800826*math.Pow(betta, 4) - 0.0499385883869711*math.Pow(betta, 3) + 0.260490073019555*math.Pow(betta, 2) - 0.638780633392638*betta + 1.33322874436111
		f = ((x-0.80)/(0.90-0.80))*(f2-f1) + f1
	case x >= 0.9 && x < 1.0:
		f2 = 0.00299986224072718*math.Pow(betta, 4) - 0.0446144019584678*math.Pow(betta, 3) + 0.250623474976823*math.Pow(betta, 2) - 0.652327954427285*betta + 1.35223617351659
		f1 = 0.00363323914800826*math.Pow(betta, 4) - 0.0499385883869711*math.Pow(betta, 3) + 0.260490073019555*math.Pow(betta, 2) - 0.638780633392638*betta + 1.333228744361110
		f = ((x-0.90)/(1.00-0.90))*(f2-f1) + f1
	case x >= 1.0 && x < 1.25:
		f1 = 0.00299986224072718*math.Pow(betta, 4) - 0.0446144019584678*math.Pow(betta, 3) + 0.250623474976823*math.Pow(betta, 2) - 0.652327954427285*betta + 1.35223617351659
		f2 = 0.0033531053559551*math.Pow(betta, 4) - 0.0486525640644669*math.Pow(betta, 3) + 0.269830627381706*math.Pow(betta, 2) - 0.71019936691933*betta + 1.39509574273625
		f = ((x-1.00)/(1.25-1.00))*(f2-f1) + f1
	case x >= 1.25 && x < 1.5:
		f2 = 0.00182411866696994*math.Pow(betta, 4) - 0.0316872884067056*math.Pow(betta, 3) + 0.208614913810006*math.Pow(betta, 2) - 0.638753657258849*betta + 1.36921673184627
		f1 = 0.0033531053559551*math.Pow(betta, 4) - 0.0486525640644669*math.Pow(betta, 3) + 0.269830627381706*math.Pow(betta, 2) - 0.710199366919334*betta + 1.3950957427362
		f = ((x-1.25)/(1.50-1.25))*(f2-f1) + f1
	case x >= 1.5 && x < 2.0:
		f1 = 0.00182411866696994*math.Pow(betta, 4) - 0.0316872884067056*math.Pow(betta, 3) + 0.208614913810006*math.Pow(betta, 2) - 0.638753657258849*betta + 1.36921673184627
		f2 = 0.00367635524977048*math.Pow(betta, 4) - 0.0521127245261192*math.Pow(betta, 3) + 0.285486147354332*math.Pow(betta, 2) - 0.770646859391461*betta + 1.44245799891866
		f = ((x-1.50)/(2.00-1.50))*(f2-f1) + f1
	default:
		f = 0.91
	}

	return f
}

// функция тупо скопирована из оригинала
func (s *FlangeService) calculateBettaV(betta, x float64) float64 {
	var f, f1, f2 float64

	switch {
	case x >= 0 && x < 0.10:
		f1 = 0.550103
		f2 = 0.0058641634223339*math.Pow(betta, 4) - 0.0748566044272087*math.Pow(betta, 3) + 0.345864798973686*math.Pow(betta, 2) - 0.708412154836844*betta + 0.980776349799449
		f = ((x-0.00)/(0.10-0.00))*(f2-f1) + f1
	case x >= 0.10 && x < 0.12:
		f1 = 0.0058641634223339*math.Pow(betta, 4) - 0.0748566044272087*math.Pow(betta, 3) + 0.345864798973686*math.Pow(betta, 2) - 0.708412154836844*betta + 0.980776349799449
		f2 = 0.00427218489206168*math.Pow(betta, 4) - 0.0576696700450847*math.Pow(betta, 3) + 0.287137267810623*math.Pow(betta, 2) - 0.646191390760794*betta + 0.961227872571484
		f = ((x-0.10)/(0.12-0.10))*(f2-f1) + f1
	case x >= 0.12 && x < 0.14:
		f1 = 0.00427218489206168*math.Pow(betta, 4) - 0.0576696700450847*math.Pow(betta, 3) + 0.287137267810623*math.Pow(betta, 2) - 0.646191390760794*betta + 0.961227872571484
		f2 = 0.00488737597016667*math.Pow(betta, 4) - 0.0658758226608894*math.Pow(betta, 3) + 0.327951390871681*math.Pow(betta, 2) - 0.738605356451387*betta + 1.02085861681089
		f = ((x-0.12)/(0.14-0.12))*(f2-f1) + f1
	case x >= 0.14 && x < 0.16:
		f1 = 0.00488737597016667*math.Pow(betta, 4) - 0.0658758226608894*math.Pow(betta, 3) + 0.327951390871681*math.Pow(betta, 2) - 0.738605356451387*betta + 1.02085861681089
		f2 = 0.00460271797088437*math.Pow(betta, 4) - 0.0648314897948872*math.Pow(betta, 3) + 0.336874456161432*math.Pow(betta, 2) - 0.786864847395865*betta + 1.05975954575613
		f = ((x-0.14)/(0.16-0.14))*(f2-f1) + f1
	case x >= 0.16 && x < 0.18:
		f1 = 0.00460271797088437*math.Pow(betta, 4) - 0.0648314897948872*math.Pow(betta, 3) + 0.336874456161432*math.Pow(betta, 2) - 0.786864847395865*betta + 1.05975954575613
		f2 = 0.00727051843506902*math.Pow(betta, 4) - 0.0961034781946584*math.Pow(betta, 3) + 0.464213008066531*math.Pow(betta, 2) - 1.00451378998363*betta + 1.1794064225988
		f = ((x-0.16)/(0.18-0.16))*(f2-f1) + f1
	case x >= 0.18 && x < 0.20:
		f1 = 0.00727051843506902*math.Pow(betta, 4) - 0.0961034781946584*math.Pow(betta, 3) + 0.464213008066531*math.Pow(betta, 2) - 1.00451378998363*betta + 1.1794064225988
		f2 = 0.00721173540592826*math.Pow(betta, 4) - 0.0969773506657811*math.Pow(betta, 3) + 0.477324201645806*math.Pow(betta, 2) - 1.05026603954099*betta + 1.21260244931179
		f = ((x-0.18)/(0.20-0.18))*(f2-f1) + f1
	case x >= 0.20 && x < 0.25:
		f1 = 0.00721173540592826*math.Pow(betta, 4) - 0.0969773506657811*math.Pow(betta, 3) + 0.477324201645806*math.Pow(betta, 2) - 1.05026603954099*betta + 1.21260244931179
		f2 = 0.00896333858044817*math.Pow(betta, 4) - 0.120259167814502*math.Pow(betta, 3) + 0.589015612657559*math.Pow(betta, 2) - 1.28321099000847*betta + 1.35440762589571
		f = ((x-0.20)/(0.25-0.20))*(f2-f1) + f1
	case x >= 0.25 && x < 0.30:
		f1 = 0.00896333858044817*math.Pow(betta, 4) - 0.120259167814502*math.Pow(betta, 3) + 0.589015612657559*math.Pow(betta, 2) - 1.28321099000847*betta + 1.35440762589571
		f2 = 0.00902846050543231*math.Pow(betta, 4) - 0.123389873098499*math.Pow(betta, 3) + 0.619559493162865*math.Pow(betta, 2) - 1.38799290494019*betta + 1.43200269598774
		f = ((x-0.25)/(0.30-0.25))*(f2-f1) + f1
	case x >= 0.30 && x < 0.35:
		f1 = 0.00902846050543231*math.Pow(betta, 4) - 0.123389873098499*math.Pow(betta, 3) + 0.619559493162865*math.Pow(betta, 2) - 1.38799290494019*betta + 1.43200269598774
		f2 = 0.0106958465978694*math.Pow(betta, 4) - 0.144294474129544*math.Pow(betta, 3) + 0.714737216488311*math.Pow(betta, 2) - 1.57772826558575*betta + 1.54552304975696
		f = ((x-0.30)/(0.35-0.30))*(f2-f1) + f1
	case x >= 0.35 && x < 0.40:
		f1 = 0.0106958465978694*math.Pow(betta, 4) - 0.144294474129544*math.Pow(betta, 3) + 0.714737216488311*math.Pow(betta, 2) - 1.57772826558575*betta + 1.54552304975696
		f2 = 0.0115878660766977*math.Pow(betta, 4) - 0.154969504488504*math.Pow(betta, 3) + 0.763360337275698*math.Pow(betta, 2) - 1.68118717343578*betta + 1.61074399314188
		f = ((x-0.35)/(0.40-0.35))*(f2-f1) + f1
	case x >= 0.40 && x < 0.45:
		f1 = 0.0115878660766977*math.Pow(betta, 4) - 0.154969504488504*math.Pow(betta, 3) + 0.763360337275698*math.Pow(betta, 2) - 1.68118717343578*betta + 1.61074399314188
		f2 = 0.0105350930422201*math.Pow(betta, 4) - 0.146169586118234*math.Pow(betta, 3) + 0.746373050400579*math.Pow(betta, 2) - 1.69254839191381*betta + 1.63032133505873
		f = ((x-0.40)/(0.45-0.40))*(f2-f1) + f1
	case x >= 0.45 && x < 0.50:
		f1 = 0.0105350930422201*math.Pow(betta, 4) - 0.146169586118234*math.Pow(betta, 3) + 0.746373050400579*math.Pow(betta, 2) - 1.69254839191381*betta + 1.63032133505873
		f2 = 0.0147272650480606*math.Pow(betta, 4) - 0.194226556641698*math.Pow(betta, 3) + 0.936899618560913*math.Pow(betta, 2) - 2.00513152244893*betta + 1.79700919379159
		f = ((x-0.45)/(0.50-0.45))*(f2-f1) + f1
	case x >= 0.50 && x < 0.60:
		f1 = 0.0147272650480606*math.Pow(betta, 4) - 0.194226556641698*math.Pow(betta, 3) + 0.936899618560913*math.Pow(betta, 2) - 2.00513152244893*betta + 1.79700919379159
		f2 = 0.0148410699644414*math.Pow(betta, 4) - 0.196084859543171*math.Pow(betta, 3) + 0.950648672240747*math.Pow(betta, 2) - 2.05089726707408*betta + 1.83003049816558
		f = ((x-0.50)/(0.60-0.50))*(f2-f1) + f1
	case x >= 0.60 && x < 0.70:
		f1 = 0.0148410699644414*math.Pow(betta, 4) - 0.196084859543171*math.Pow(betta, 3) + 0.950648672240747*math.Pow(betta, 2) - 2.05089726707408*betta + 1.83003049816558
		f2 = 0.0138984227474326*math.Pow(betta, 4) - 0.188197118526749*math.Pow(betta, 3) + 0.935272559601286*math.Pow(betta, 2) - 2.0628061986072*betta + 1.85045077372066
		f = ((x-0.60)/(0.70-0.60))*(f2-f1) + f1
	case x >= 0.7 && x < 0.8:
		f1 = 0.0138984227474326*math.Pow(betta, 4) - 0.188197118526749*math.Pow(betta, 3) + 0.935272559601286*math.Pow(betta, 2) - 2.0628061986072*betta + 1.85045077372066
		f2 = 0.0161105067667355*math.Pow(betta, 4) - 0.214306039151065*math.Pow(betta, 3) + 1.04396963732146*math.Pow(betta, 2) - 2.2533121521015*betta + 1.95659929065701
		f = ((x-0.70)/(0.80-0.70))*(f2-f1) + f1
	case x >= 0.8 && x < 0.9:
		f1 = 0.0161105067667355*math.Pow(betta, 4) - 0.214306039151065*math.Pow(betta, 3) + 1.04396963732146*math.Pow(betta, 2) - 2.2533121521015*betta + 1.95659929065701
		f2 = 0.0171585631685489*math.Pow(betta, 4) - 0.225971160127992*math.Pow(betta, 3) + 1.09015746815975*math.Pow(betta, 2) - 2.33460469421229*betta + 2.0015405396793
		f = ((x-0.80)/(0.90-0.80))*(f2-f1) + f1
	case x >= 0.9 && x < 1:
		f1 = 0.0171585631685489*math.Pow(betta, 4) - 0.225971160127992*math.Pow(betta, 3) + 1.09015746815975*math.Pow(betta, 2) - 2.33460469421229*betta + 2.0015405396793
		f2 = 0.0181832719615506*math.Pow(betta, 4) - 0.238474483719479*math.Pow(betta, 3) + 1.14422460739396*math.Pow(betta, 2) - 2.43362479844938*betta + 2.0574674694669
		f = ((x-0.90)/(1.00-0.90))*(f2-f1) + f1
	case x >= 1 && x < 1.25:
		f1 = 0.0181832719615506*math.Pow(betta, 4) - 0.238474483719479*math.Pow(betta, 3) + 1.14422460739396*math.Pow(betta, 2) - 2.43362479844938*betta + 2.0574674694669
		f2 = 0.0193245080393135*math.Pow(betta, 4) - 0.252719162624534*math.Pow(betta, 3) + 1.20730036223498*math.Pow(betta, 2) - 2.55097264670514*betta + 2.12464092473434
		f = ((x-1.00)/(1.25-1.00))*(f2-f1) + f1
	case x >= 1.25 && x < 1.50:
		f1 = 0.0193245080393135*math.Pow(betta, 4) - 0.252719162624534*math.Pow(betta, 3) + 1.20730036223498*math.Pow(betta, 2) - 2.55097264670514*betta + 2.12464092473434
		f2 = 0.0210833381428893*math.Pow(betta, 4) - 0.275706486072612*math.Pow(betta, 3) + 1.31244377792127*math.Pow(betta, 2) - 2.74660578969882*betta + 2.23635760844775
		f = ((x-1.25)/(1.50-1.25))*(f2-f1) + f1
	case x >= 1.5 && x <= 2.00:
		f1 = 0.0210833381428893*math.Pow(betta, 4) - 0.275706486072612*math.Pow(betta, 3) + 1.31244377792127*math.Pow(betta, 2) - 2.74660578969882*betta + 2.23635760844775
		f2 = 0.021767904764524*math.Pow(betta, 4) - 0.285402818872975*math.Pow(betta, 3) + 1.36155533631857*math.Pow(betta, 2) - 2.84937881931514*betta + 2.29865470248542
		f = ((x-1.5)/(2.00-1.5))*(f2-f1) + f1
	default:
		f = 0.55
	}

	return f
}

func (s *FlangeService) calculateF(betta, x float64) float64 {
	var f, f1, f2 float64

	switch {
	case x >= 0 && x < 0.05:
		f1 = ((betta-1.0)/(5.0-1.0))*(25.0-1.0) + 1.0
		f2 = ((betta-1.03333)/(5.0-1.03333))*(23.33333-1.0) + 1.0
		f = ((x-0.0)/(0.05-0.0))*(f2-f1) + f1
	case x >= 0.05 && x < 0.10:
		f1 = ((betta-1.03333)/(5.0-1.03333))*(23.33333-1.0) + 1.0
		f2 = ((betta-1.12)/(5.0-1.12))*(21.0-1.0) + 1.0
		f = ((x-0.05)/(0.10-0.05))*(f2-f1) + f1
	case x >= 0.10 && x < 0.15:
		f1 = ((betta-1.2)/(5.0-1.2))*(21.0-1.0) + 1.0
		f2 = ((betta-1.18)/(5.0-1.18))*(19.0-1.0) + 1.0
		f = ((x-0.10)/(0.15-0.10))*(f2-f1) + f1
	case x >= 0.15 && x < 0.20:
		f1 = ((betta-1.18)/(5.0-1.18))*(19.0-1.0) + 1.0
		f2 = ((betta-1.24)/(5.0-1.24))*(16.2-1.0) + 1.0
		f = ((x-0.15)/(0.20-0.15))*(f2-f1) + f1
	case x >= 0.20 && x < 0.25:
		f1 = ((betta-1.24)/(5.0-1.24))*(16.2-1.0) + 1.0
		f2 = ((betta-1.33)/(5.0-1.33))*(15.0-1.0) + 1.0
		f = ((x-0.20)/(0.25-0.20))*(f2-f1) + f1
	case x >= 0.25 && x < 0.30:
		f1 = ((betta-1.33)/(5.0-1.33))*(15.0-1.0) + 1.0
		f2 = ((betta-1.4)/(5.0-1.4))*(13.5-1.0) + 1.0
		f = ((x-0.25)/(0.30-0.25))*(f2-f1) + f1
	case x >= 0.30 && x < 0.35:
		f1 = ((betta-1.4)/(5.0-1.4))*(13.5-1.0) + 1.0
		f2 = ((betta-1.48)/(5.0-1.48))*(12.5-1.0) + 1.0
		f = ((x-0.30)/(0.35-0.30))*(f2-f1) + f1
	case x >= 0.35 && x < 0.40:
		f1 = ((betta-1.48)/(5.0-1.48))*(12.5-1.0) + 1.0
		f2 = ((betta-1.6)/(5.0-1.6))*(10.5-1.0) + 1.0
		f = ((x-0.35)/(0.40-0.35))*(f2-f1) + f1
	case x >= 0.40 && x < 0.45:
		f1 = ((betta-1.6)/(5.0-1.6))*(10.5-1.0) + 1.0
		f2 = ((betta-1.7)/(5.0-1.7))*(9.0-1.0) + 1.0
		f = ((x-0.40)/(0.45-0.40))*(f2-f1) + f1
	case x >= 0.45 && x < 0.50:
		f1 = ((betta-1.7)/(5.0-1.7))*(9.0-1.0) + 1.0
		f2 = ((betta-1.8)/(5.0-1.8))*(8.25-1.0) + 1.0
		f = ((x-0.45)/(0.50-0.45))*(f2-f1) + f1
	case x >= 0.50 && x < 0.60:
		f1 = ((betta-1.8)/(5.0-1.8))*(8.25-1.0) + 1.0
		f2 = ((betta-2.05)/(5.0-2.05))*(6.5-1.0) + 1.0
		f = ((x-0.50)/(0.60-0.50))*(f2-f1) + f1
	case x >= 0.60 && x < 0.70:
		f1 = ((betta-2.05)/(5.0-2.05))*(6.5-1.0) + 1.0
		f2 = ((betta-2.35)/(5.0-2.35))*(5.0-1.0) + 1.0
		f = ((x-0.60)/(0.70-0.60))*(f2-f1) + f1
	case x >= 0.70 && x < 0.80:
		f1 = ((betta-2.35)/(5.0-2.35))*(5.0-1.0) + 1.0
		f2 = ((betta-2.65)/(5.0-2.65))*(4.42-1.0) + 1.0
		f = ((x-0.70)/(0.80-0.70))*(f2-f1) + f1
	case x >= 0.80 && x < 0.90:
		f1 = ((betta-2.65)/(5.0-2.65))*(4.42-1.0) + 1.0
		f2 = ((betta-3.0)/(5.0-3.0))*(3.5-1.0) + 1.0
		f = ((x-0.80)/(0.90-0.80))*(f2-f1) + f1
	case x >= 0.90 && x < 1.00:
		f1 = ((betta-3.0)/(5.0-3.0))*(3.5-1.0) + 1.0
		f2 = ((betta-3.4)/(5.0-3.4))*(2.7-1.0) + 1.0
		f = ((x-0.90)/(1.00-0.90))*(f2-f1) + f1
	case x >= 1.00 && x < 1.10:
		f1 = ((betta-3.4)/(5.0-3.4))*(2.7-1.0) + 1.0
		f2 = ((betta-3.8)/(5.0-3.8))*(2.4-1.0) + 1.0
		f = ((x-1.00)/(1.10-1.00))*(f2-f1) + f1
	case x >= 1.10 && x < 1.20:
		f1 = ((betta-3.8)/(5.0-3.8))*(2.4-1.0) + 1.0
		f2 = ((betta-4.2)/(5.0-4.2))*(1.9-1.0) + 1.0
		f = ((x-1.10)/(1.20-1.10))*(f2-f1) + f1
	case x >= 1.20 && x <= 1.30:
		f1 = ((betta-4.2)/(5.0-4.2))*(1.9-1.0) + 1.0
		f2 = ((betta-4.63)/(5.0-4.63))*(1.3-1.0) + 1.0
		f = ((x-1.20)/(1.30-1.20))*(f2-f1) + f1
	default:
		f = 1.0
	}

	return f
}
