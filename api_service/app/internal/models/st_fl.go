package models

type StFlDTO struct {
	StandId   string `json:"standId" binding:"required"`
	FlangeIds string `json:"flIds" binding:"required"`
}
