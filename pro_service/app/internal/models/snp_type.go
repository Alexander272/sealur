package models

type SNPType struct {
	Id    string `db:"id"`
	Title string `db:"title"`
	Code  string `db:"code"`
}

type SNPTypeWithFlange struct {
	Id          string `db:"id"`
	Title       string `db:"title"`
	Code        string `db:"code"`
	TypeId      string `db:"type_id"`
	TypeTitle   string `db:"type_title"`
	TypeCode    string `db:"type_code"`
	Description string `db:"description"`
	HasD4       bool   `db:"has_d4"`
	HasD3       bool   `db:"has_d3"`
	HasD2       bool   `db:"has_d2"`
	HasD1       bool   `db:"has_d1"`
}
