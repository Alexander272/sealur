package models

type FlangeDTO struct {
	Title string `json:"title" binding:"required"`
	Short string `json:"short" binding:"required"`
}
