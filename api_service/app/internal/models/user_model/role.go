package user_model

type UserRole struct {
	UserId  string `json:"userId"`
	Service string `json:"service" binding:"required"`
	Role    string `json:"role" binding:"required"`
}
