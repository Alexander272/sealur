package models

// TODO надо назвать поля по другому
type FlangeSize struct {
	D  float64 `db:"D"`
	Pn float64 `db:"pn"`
	D1 float64 `db:"D1"`
	D2 float64 `db:"D2"`
	// D6       float32 `db:"D6"`
	// D7       float32 `db:"D7"`
	// H        float32 `db:"h"`
	S0       float64 `db:"S0"`
	S1       float64 `db:"S1"`
	B        float64 `db:"b"`
	Lenght   float64 `db:"lenght"`
	Count    int32   `db:"count"`
	Diameter int32   `db:"diameter"`
	Area     float64 `db:"area"`
}

type InitialDataFlange struct {
	DOut         float64
	D            float64
	Dk           float64
	Dnk          float64
	Ds           float64
	H            float64
	Hk           float64
	S0           float64
	S1           float64
	L            float64
	D6           float64
	C            float64
	Tf           float64
	Tk           float64
	AlphaK       float64
	EpsilonKAt20 float64
	EpsilonK     float64
	SigmaKAt20   float64
	SigmaK       float64
	AlphaF       float64
	EpsilonAt20  float64
	Epsilon      float64
	SigmaAt20    float64
	Sigma        float64
	SigmaM       float64
	SigmaMAt20   float64
	SigmaR       float64
	SigmaRAt20   float64
	Material     string
	Count        int32
	Diameter     int32
	Area         float64
}

type CalculatedData struct {
	B      float64
	A      float64
	Ds     float64
	E      float64
	Se     float64
	Zak    float64
	Betta  float64
	X      float64
	L0     float64
	K      float64
	BettaT float64
	BettaU float64
	BettaY float64
	BettaZ float64
	BettaF float64
	BettaV float64
	F      float64
	Lymda  float64
	Yf     float64
	Yk     float64
	Psik   float64
	Yfn    float64
	Yfc    float64
}

type CalculatedBasis struct {
	Qt      float64
	Pb1     float64
	Pb2     float64
	Pbm     float64
	Pbr     float64
	SigmaB1 float64
	SigmaB2 float64
	DSigmaM float64
	DSigmaR float64
	Q       float64
	Mkp     float64
	Mkp1    float64
	Mrek    float64
	Qrek    float64
	Mmax    float64
	Qmax    float64
}

type CalculatedStrength struct {
	Mkp       float64
	Mkp1      float64
	Cf        float64
	Dzv       float64
	MM        float64
	MMk       float64
	Mpk       float64
	Mp        float64
	SigmaM1   float64
	SigmaM0   float64
	SigmaT    float64
	SigmaR    float64
	SigmaTp   float64
	SigmaRp   float64
	SigmaK    float64
	SigmaP1   float64
	SigmaP0   float64
	SigmaMp   float64
	SigmaMpm  float64
	SigmaMp0  float64
	SigmaMpm0 float64
	SigmaMop  float64
	SigmaKp   float64

	Teta   float64
	DTeta  float64
	DTetaK float64
	TetaK  float64

	Max1, Max2, Max3, Max4, Max5, Max6, Max7, Max8, Max9, Max10, Max11 float64
}
