package models

type OrderDTO struct {
	Id     string `json:"id"`
	Count  int32  `json:"count"`
	UserId string `json:"userId" binding:"required"`
}

type PositionDTO struct {
	Designation string `json:"designation"`
	Descriprion string `json:"description"`
	Count       string `json:"count"`
	Sizes       string `json:"sizes"`
	Drawing     string `json:"drawing"`
}
