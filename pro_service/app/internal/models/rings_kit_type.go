package models

type RingsKitType struct {
	Id          string `db:"id"`
	Code        string `db:"code"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Image       string `db:"image"`
	Designation string `db:"designation"`
}
