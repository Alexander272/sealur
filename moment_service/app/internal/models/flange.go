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
	Length   float64 `db:"length"`
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
	Length  float64 `db:"length"`
	Count   int32   `db:"count"`
	BoltId  string  `db:"bolt_id"`
}

type TypeFlangeDTO struct {
	Id    string `db:"id"`
	Title string `db:"title"`
	Label string `db:"label"`
}

type GetBasisSize struct {
	IsUseRow bool
	StandId  string `db:"stand_id"`
	Row      string `db:"row"`
}

type Data struct {
	Work       string
	Flanges    string
	SameFlange string
	Embedded   string
	Type       string
	Condition  string
}
