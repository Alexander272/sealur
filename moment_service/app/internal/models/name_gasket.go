package models

type NameGasketDTO struct {
	Id     string `db:"id"`
	FinId  string `db:"fin_id"`
	NumId  string `db:"num_id"`
	PresId string `db:"pres_id"`
	Title  string `db:"title"`
}

type FullNameGasketDTO struct {
	Id        string  `db:"id"`
	FinId     string  `db:"fin_id"`
	NumId     string  `db:"num_id"`
	PresId    string  `db:"pres_id"`
	Title     string  `db:"title"`
	SizeLong  float64 `db:"size_long"`
	SizeTrans float64 `db:"size_trans"`
	Width     float64 `db:"width"`
	Thick1    float64 `db:"thick1"`
	Thick2    float64 `db:"thick2"`
	Thick3    float64 `db:"thick3"`
	Thick4    float64 `db:"thick4"`
}

type NameGasketSizeDTO struct {
	Id        string  `db:"id"`
	SizeLong  float64 `db:"size_long"`
	SizeTrans float64 `db:"size_trans"`
	Width     float64 `db:"width"`
	Thick1    float64 `db:"thick1"`
	Thick2    float64 `db:"thick2"`
	Thick3    float64 `db:"thick3"`
	Thick4    float64 `db:"thick4"`
}
