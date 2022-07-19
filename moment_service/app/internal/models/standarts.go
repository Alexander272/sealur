package models

type StandartDTO struct {
	Id     string `db:"id"`
	Title  string `db:"title"`
	TypeId string `db:"type_id"`
}
