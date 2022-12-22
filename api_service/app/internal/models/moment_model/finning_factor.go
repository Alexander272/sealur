package moment_model

type FinningFactorDTO struct {
	DevId string `json:"devId" binding:"required"`
	Value string `json:"value" binding:"required"`
}
