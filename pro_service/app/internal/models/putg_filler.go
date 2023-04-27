package models

type PutgFiller struct {
	Id          string `db:"id"`
	Temperature string `db:"temperature"`
	Title       string `db:"title"`
	Code        string `db:"code"`
	Description string `db:"description"`
	Designation string `db:"designation"`
}
