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
}

type InitialDataFlange struct {
	DOut        float64
	D           float64
	H           float64
	S0          float64
	S1          float64
	L           float64
	D6          float64
	C           float64
	Tf          float64
	AlphaF      float64
	EpsilonAt20 float64
	Epsilon     float64
	SigmaAt20   float64
	Sigma       float64
	SigmaM      float64
	SigmaMAt20  float64
	SigmaR      float64
	SigmaRAt20  float64
	Material    string
}

type InitialDataBolt struct {
}
