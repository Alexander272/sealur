package models

type Role struct {
	Id      string `db:"id"`
	UserId  string `db:"user_id"`
	Service string `db:"service"`
	Role    string `db:"role"`
}
