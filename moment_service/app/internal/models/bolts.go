package models

type BoltSize struct {
	Diameter float64 `db:"diameter"`
	Area     float64 `db:"area"`
}

type BoltsDTO struct {
	Id       string  `db:"id"`
	Title    string  `db:"title"`
	Diameter float64 `db:"diameter"`
	Area     float64 `db:"area"`
	IsInch   bool    `db:"is_inch"`
}
