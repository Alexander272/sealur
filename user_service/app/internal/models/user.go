package models

type User struct {
	Id           string `db:"id"`
	Organization string `db:"organization"`
	Name         string `db:"name"`
	Email        string `db:"email"`
	City         string `db:"city"`
	Position     string `db:"position"`
	Phone        string `db:"phone"`
	Password     string `db:"password"`
}
