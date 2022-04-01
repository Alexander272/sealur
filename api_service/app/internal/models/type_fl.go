package models

type TypeFlDTO struct {
	Title string `json:"title" binding:"required"`
	Descr string `json:"descr"`
	Short string `json:"short"`
	Basis bool   `json:"basis"`
}
