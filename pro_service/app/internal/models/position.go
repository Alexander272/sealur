package models

type Position struct {
	Id          string `db:"id"`
	Designation string `db:"designation"`
	Descriprion string `db:"description"`
	Count       int32  `db:"count"`
	Sizes       string `db:"sizes"`
	Drawing     string `db:"drawing"`
	OrderId     string `db:"order_id"`
}

// id, title, amount, type, count
type PositionNew struct {
	Id     string `db:"id"`
	Title  string `db:"title"`
	Amount string `db:"amount"`
	Type   string `db:"type"`
	Count  int64  `db:"count"`
}
