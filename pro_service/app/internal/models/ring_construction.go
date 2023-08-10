package models

type RingConstruction struct {
	Id          string `db:"id"`
	TypeId      string `db:"type_id"`
	Code        string `db:"code"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Image       string `db:"image"`
}
