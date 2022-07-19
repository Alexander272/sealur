package models

type BoltsDTO struct {
	Id       string  `db:"id"`
	Title    string  `db:"title"`
	Diameter int32   `db:"diameter"`
	Area     float64 `db:"area"`
}
