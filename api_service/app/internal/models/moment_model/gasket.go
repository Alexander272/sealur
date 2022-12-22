package moment_model

type GasketDTO struct {
	Title string `json:"title" binding:"required"`
}

type TypeGasketDTO struct {
	Title string `json:"title"`
	Label string `json:"label"`
}

type GasketDataDTO struct {
	GasketId        string  `json:"gasketId"`
	PermissiblePres float64 `json:"permissiblePres"`
	Compression     float64 `json:"compression"`
	Epsilon         float64 `json:"epsilon"`
	Thickness       float64 `json:"thickness"`
	TypeId          string  `json:"typeId"`
}

type ManyGasketDataDTO struct {
	GasketId string                   `json:"gasketId"`
	TypeId   string                   `json:"typeId"`
	Data     []ManyGasketDataDTO_Data `json:"data"`
}

type ManyGasketDataDTO_Data struct {
	PermissiblePres float64 `json:"permissiblePres"`
	Compression     float64 `json:"compression"`
	Epsilon         float64 `json:"epsilon"`
	Thickness       float64 `json:"thickness"`
}
