package models

type PutgmImageDTO struct {
	Form   string `json:"form" binding:"required"`
	Gasket string `json:"gasket" binding:"required"`
	Url    string `json:"url" binding:"required"`
}
