package models

type PutgmImageDTO struct {
	Form   string `json:"form"  form:"form"  binding:"required"`
	Gasket string `json:"gasket" form:"gasket" binding:"required"`
	Url    string `json:"url" form:"url" binding:"required"`
}
