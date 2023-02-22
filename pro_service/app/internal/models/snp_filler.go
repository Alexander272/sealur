package models

type SNPFiller struct {
	Id          string `db:"id"`
	Title       string `db:"title"`
	Code        string `db:"code"`
	Description string `db:"description"`
}
