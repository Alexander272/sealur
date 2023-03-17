package models

// type User struct {
// 	Id           string `db:"id"`
// 	Organization string `db:"organization"`
// 	Name         string `db:"name"`
// 	Email        string `db:"email"`
// 	City         string `db:"city"`
// 	Position     string `db:"position"`
// 	Phone        string `db:"phone"`
// 	Password     string `db:"password"`
// 	Login        string `db:"login"`
// 	Count        int    `db:"count"`
// }

type User struct {
	Id       string `db:"id"`
	Company  string `db:"company"`
	Address  string `db:"address"`
	Inn      string `db:"inn"`
	Kpp      string `db:"kpp"`
	Region   string `db:"region"`
	City     string `db:"city"`
	Name     string `db:"name"`
	Position string `db:"position"`
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	RoleCode string `db:"role_code"`
	Password string `db:"password"`
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
