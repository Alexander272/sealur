package models

type RingModifying struct {
	Id          string `db:"id"`
	Code        string `db:"code"`
	Description string `db:"description"`
}
