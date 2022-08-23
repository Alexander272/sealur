package moment_model

type BoltDTO struct {
	Title    string  `json:"title"`
	Diameter float64 `json:"diameter"`
	Area     float64 `json:"area"`
	IsInch   bool    `json:"isInch"`
}
