package models

type StandDTO struct {
	Title string `json:"title" binding:"required"`
}
