package models

type PutgType struct {
	Id           string  `db:"id"`
	Title        string  `db:"title"`
	Code         string  `db:"code"`
	FillerId     string  `db:"filler_id"`
	MinThickness float64 `db:"min_thickness"`
	MaxThickness float64 `db:"max_thickness"`
	Description  string  `db:"description"`
	TypeCode     string  `db:"type_code"`
}
