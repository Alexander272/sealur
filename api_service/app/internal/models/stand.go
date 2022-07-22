package models

type StandDTO struct {
	Title string `json:"title" binding:"required"`
}

type MomentStandartDTO struct {
	Title  string `json:"title"`
	TypeId string `json:"typeId"`
}
