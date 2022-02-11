package models

type TypeFlDTO struct {
	Title string `json:"title" binding:"required"`
	Desc  string `json:"desc" binding:"required"`
	Short string `json:"short" binding:"required"`
	Basis string `json:"basis" binding:"required"`
}
