package models

type SectionExecutionDTO struct {
	Id    string `db:"id"`
	DevId string `db:"dev_id"`
	Value string `db:"value"`
}
