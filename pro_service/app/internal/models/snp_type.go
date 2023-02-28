package models

type SNPType struct {
	Id    string `db:"id"`
	Title string `db:"title"`
	Code  string `db:"code"`
}

type SNPTypeWithFlange struct {
	Id        string `db:"id"`
	Title     string `db:"title"`
	Code      string `db:"code"`
	TypeId    string `db:"type_id"`
	TypeTitle string `db:"type_title"`
	TypeCode  string `db:"type_code"`
}
