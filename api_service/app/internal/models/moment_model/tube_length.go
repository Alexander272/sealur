package moment_model

type TubeLengthDTO struct {
	DevId string `json:"devId" binding:"required"`
	Value string `json:"value" binding:"required"`
}
