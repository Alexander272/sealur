package user_model

type SignUp struct {
	Organization string `json:"organization" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"email,required"`
	City         string `json:"city"`
	Position     string `json:"position"`
	Phone        string `json:"phone"`
}

type SignIn struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ConfirmUser struct {
	Id       string     `json:"id" binding:"required"`
	Login    string     `json:"login" binding:"required"`
	Password string     `json:"password" binding:"required"`
	Roles    []UserRole `json:"roles" binding:"required"`
}

type UpdateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
