package moment_model

type GasketDTO struct {
	Title string `json:"title" binding:"required"`
}

type TypeGasketDTO struct {
	Title string `json:"title" binding:"required"`
}

type GasketDataDTO struct {
	GasketId        string  `json:"gasketId"`
	PermissiblePres float64 `json:"permissiblePres"`
	Compression     float64 `json:"compression"`
	Epsilon         float64 `json:"epsilon"`
	Thickness       float64 `json:"thickness"`
	TypeId          string  `json:"typeId"`
}
