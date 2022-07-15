package constants

const (
	IsolatedFlatTb float64 = 0.97
	IsolatedFlatTf float64 = 1
	// NonIsolatedFlatTb float64 = 0.95
	NonIsolatedFlatTb float64 = 0.85
	NonIsolatedFlatTf float64 = 0.96

	IsolatedFreeTb    float64 = 0.9
	IsolatedFreeTk    float64 = 0.97
	IsolatedFreeTf    float64 = 1
	NonIsolatedFreeTb float64 = 0.81
	NonIsolatedFreeTk float64 = 0.9
	NonIsolatedFreeTf float64 = 0.96

	SigmaM float64 = 1.5
	SigmaR float64 = 3

	B0    float64 = 3.8
	Bp    float64 = 15
	BoltD float64 = 0.28
	PinD  float64 = 0.56

	WorkKyp float64 = 1
	TestKyp float64 = 1.35

	UncontrollableKyz  float64 = 1
	ControllableKyz    float64 = 1.1
	ControllablePinKyz float64 = 1.3

	NoLoadKyt float64 = 1
	LoadKyt   float64 = 1.3

	MinDiameter int32   = 20
	MaxDiameter int32   = 52
	MaxSigmaB   float64 = 120
)
