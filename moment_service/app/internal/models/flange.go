package models

type FlangeSize struct {
	Id       string  `db:"id"`
	Pn       float64 `db:"pn"`
	D        float64 `db:"d"`
	D6       float64 `db:"d6"`
	DOut     float64 `db:"d_out"`
	H        float64 `db:"h"`
	S0       float64 `db:"s0"`
	S1       float64 `db:"s1"`
	Lenght   float64 `db:"lenght"`
	Count    int32   `db:"count"`
	Diameter int32   `db:"diameter"`
	Area     float64 `db:"area"`
}

type FlangeSizeDTO struct {
	Id      string  `db:"id"`
	StandId string  `db:"stand_id"`
	Pn      float64 `db:"pn"`
	D       float64 `db:"d"`
	D6      float64 `db:"d6"`
	DOut    float64 `db:"d_out"`
	H       float64 `db:"h"`
	S0      float64 `db:"s0"`
	S1      float64 `db:"s1"`
	Lenght  float64 `db:"lenght"`
	Count   int32   `db:"count"`
	BoltId  string  `db:"bolt_id"`
}

type TypeFlangeDTO struct {
	Id    string `db:"id"`
	Title string `db:"title"`
}
