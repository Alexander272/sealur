package moment_model

type NameGasketDTO struct {
	FinId     string  `json:"finId" binding:"required"`
	NumId     string  `json:"numId" binding:"required"`
	PresId    string  `json:"presId" binding:"required"`
	Title     string  `json:"title" binding:"required"`
	SizeLong  float64 `json:"sizeLong"`
	SizeTrans float64 `json:"sizeTrans"`
	Width     float64 `json:"width"`
	Thick1    float64 `json:"thick1"`
	Thick2    float64 `json:"thick2"`
	Thick3    float64 `json:"thick3"`
	Thick4    float64 `json:"thick4"`
}
