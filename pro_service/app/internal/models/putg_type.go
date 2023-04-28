package models

type PutgType struct {
	Id       string `db:"id"`
	Title    string `db:"title"`
	Code     string `db:"code"`
	FillerId string `db:"filler_id"`
}
