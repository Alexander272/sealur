package models

type NumberOfMovesDTO struct {
	Id      string `db:"id"`
	DevId   string `db:"dev_id"`
	CountId string `db:"count_id"`
	Value   string `db:"value"`
}
