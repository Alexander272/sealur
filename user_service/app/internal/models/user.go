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
	Login        string `db:"login"`
}

type ConfirmUser struct {
	Name  string `db:"name"`
	Email string `db:"email"`
}

type DeleteUser struct {
	Name  string `db:"name"`
	Email string `db:"email"`
}

type Count struct {
	Count int32 `db:"count"`
}
