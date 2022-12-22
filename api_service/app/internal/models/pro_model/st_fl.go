package pro_model

type StFlDTO struct {
	StandId  string `json:"standId" binding:"required"`
	FlangeId string `json:"flangeId" binding:"required"`
}
