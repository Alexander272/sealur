package models

type FinningFactorDTO struct {
	Id    string `db:"id"`
	DevId string `db:"dev_id"`
	Value string `db:"value"`
}
