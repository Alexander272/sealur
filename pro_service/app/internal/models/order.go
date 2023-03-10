package models

type Order struct {
	Id    string `db:"id"`
	Date  string `db:"date"`
	Count int32  `db:"count_position"`
}

type OrderNew struct {
	Id     string `db:"id"`
	Date   string `db:"date"`
	Count  int64  `db:"count_position"`
	Number int64  `db:"number"`
}

type OrderWithPosition struct {
	Id            string `db:"id"`
	Date          string `db:"date"`
	Count         int64  `db:"count_position"`
	Number        int64  `db:"number"`
	PositionId    string `db:"position_id"`
	Title         string `db:"title"`
	Amount        string `db:"amount"`
	PositionCount int64  `db:"position_count"`
}
