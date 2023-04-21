package models

type OrderNew struct {
	Id     string `db:"id"`
	Date   string `db:"date"`
	Count  int64  `db:"count_position"`
	Number int64  `db:"number"`
	Status string `db:"status"`
	Info   string `db:"info"`
}

type ManagerOrder struct {
	Id      string `db:"id"`
	Date    string `db:"date"`
	UserId  string `db:"user_id"`
	Company string `db:"company"`
	Count   int64  `db:"count_position"`
	Number  int64  `db:"number"`
	Status  string `db:"status"`
	Info    string `db:"info"`
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
	Info          string `db:"info"`
}
