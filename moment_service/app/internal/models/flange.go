package models

type FlangeSize struct {
	Id       string  `db:"id"`
	Pn       float64 `db:"pn"`
	D        float64 `db:"d"`
	Dn       string  `db:"dn"`
	D6       float64 `db:"d6"`
	DOut     float64 `db:"d_out"`
	X        float64 `db:"x"`
	A        float64 `db:"a"`
	H        float64 `db:"h"`
	S0       float64 `db:"s0"`
	S1       float64 `db:"s1"`
	Length   float64 `db:"length"`
	Count    int32   `db:"count"`
	Diameter float64 `db:"diameter"`
	Area     float64 `db:"area"`
	IsEmptyD bool    `db:"is_empty_d"`
}

type FlangeSizeDTO struct {
	Id      string  `db:"id"`
	StandId string  `db:"stand_id"`
	Pn      float64 `db:"pn"`
	Dn      string  `db:"dn"`
	Dmm     float64 `db:"dmm"`
	D       float64 `db:"d"`
	D6      float64 `db:"d6"`
	DOut    float64 `db:"d_out"`
	X       float64 `db:"x"`
	A       float64 `db:"a"`
	H       float64 `db:"h"`
	S0      float64 `db:"s0"`
	S1      float64 `db:"s1"`
	Length  float64 `db:"length"`
	Count   int32   `db:"count"`
	BoltId  string  `db:"bolt_id"`
	Row     int32   `db:"row"`
}

type TypeFlangeDTO struct {
	Id    string `db:"id"`
	Title string `db:"title"`
	Label string `db:"label"`
}

type GetBasisSize struct {
	IsUseRow bool
	StandId  string `db:"stand_id"`
	Row      int32  `db:"row"`
	IsInch   bool   `db:"is_inch"`
}

type Data struct {
	Work       string
	Flanges    string
	SameFlange string
	Embedded   string
	Type       string
	Condition  string
}
