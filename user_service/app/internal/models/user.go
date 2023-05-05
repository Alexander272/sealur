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
	Id        string `db:"id"`
	Company   string `db:"company"`
	Address   string `db:"address"`
	Inn       string `db:"inn"`
	Kpp       string `db:"kpp"`
	Region    string `db:"region"`
	City      string `db:"city"`
	Name      string `db:"name"`
	Position  string `db:"position"`
	Email     string `db:"email"`
	Phone     string `db:"phone"`
	RoleCode  string `db:"role_code"`
	Password  string `db:"password"`
	ManagerId string `db:"manager_id"`
	Confirmed bool   `db:"confirmed"`
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

type Analytics struct {
	UserCount         int64 `db:"user_count"`
	LinkCount         int64 `db:"link_count"`
	RegisterCount     int64 `db:"register_count"`
	RegisterLinkCount int64 `db:"register_link_count"`
}

type UserAnalytics struct {
	Id        string `db:"id"`
	Company   string `db:"company"`
	Position  string `db:"position"`
	Phone     string `db:"phone"`
	Email     string `db:"email"`
	Date      string `db:"date"`
	Name      string `db:"name"`
	ManagerId string `db:"manager_id"`
	UseLink   bool   `db:"use_link"`
	Manager   string `db:"manager"`
	HasOrder  bool   `db:"has_order"`
}
