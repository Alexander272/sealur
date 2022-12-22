package moment_model

type EnvDTO struct {
	Title string `json:"title" binding:"required"`
}

type EnvDataDTO struct {
	EnvId        string  `json:"envId"`
	GasketId     string  `json:"gasketId"`
	M            float64 `json:"m"`
	SpecificPres float64 `json:"specificPres"`
}

type ManyEnvDataDTO struct {
	GasketId string                `json:"gasketId"`
	Data     []ManyEnvDataDTO_Data `json:"data"`
}

type ManyEnvDataDTO_Data struct {
	EnvId        string  `json:"envId"`
	M            float64 `json:"m"`
	SpecificPres float64 `json:"specificPres"`
}
