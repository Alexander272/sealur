package models

type PressureDTO struct {
	Id    string  `db:"id"`
	Value float64 `db:"value"`
}
