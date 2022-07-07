package models

type FlangeSize struct {
	D  float32 `db:"D"`
	Pn float32 `db:"pn"`
	D1 float32 `db:"D1"`
	D2 float32 `db:"D2"`
	// D6       float32 `db:"D6"`
	// D7       float32 `db:"D7"`
	S0 float32 `db:"S0"`
	S1 float32 `db:"S1"`
	// H        float32 `db:"h"`
	B        float32 `db:"b"`
	Lenght   float32 `db:"lenght"`
	Count    int32   `db:"count"`
	Diameter string  `db:"diameter"`
}
