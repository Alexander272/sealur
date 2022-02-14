package models

type StFlDTO struct {
	StandId  string `json:"standId" binding:"required"`
	FlangeId string `json:"flangeId" binding:"required"`
}
