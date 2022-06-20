package models

type Role struct {
	Id      string `db:"id"`
	Service string `db:"service"`
	Role    string `db:"role"`
}
