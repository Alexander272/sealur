package models

type Order struct {
	Id    string `db:"id"`
	Date  string `db:"date"`
	Count int32  `db:"count_position"`
}

type Position struct {
	Id          string `db:"id"`
	Designation string `db:"designation"`
	Descriprion string `db:"description"`
	Count       int32  `db:"count"`
	Sizes       string `db:"sizes"`
	Drawing     string `db:"drawing"`
	OrderId     string `db:"order_id"`
}
