package models

type Ip struct {
	Ip     string `db:"ip"`
	Date   string `db:"date"`
	UserId string `db:"user_id"`
}
