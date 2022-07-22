package models

type EnvDTO struct {
	Title string `json:"title" binding:"required"`
}

type EnvDataDTO struct {
	EnvId        string  `json:"envId"`
	GasketId     string  `json:"gasketId"`
	M            float64 `json:"m"`
	SpecificPres float64 `json:"specificPres"`
}
