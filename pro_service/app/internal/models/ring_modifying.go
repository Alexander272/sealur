package models

type RingModifying struct {
	Id          string `db:"id"`
	Code        string `db:"code"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Designation string `db:"designation"`
}
