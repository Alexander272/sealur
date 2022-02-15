package models

type TypeFlDTO struct {
	Title string `json:"title" binding:"required"`
	Desc  string `json:"desc"`
	Short string `json:"short"`
	Basis bool   `json:"basis"`
}
