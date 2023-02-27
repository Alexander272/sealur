package models

type SNPFiller struct {
	Id           string `db:"id"`
	Title        string `db:"title"`
	AnotherTitle string `db:"another_title"`
	Code         string `db:"code"`
	Description  string `db:"description"`
	Designation  string `db:"designation"`
	Temperature  string `db:"temperature"`
}
