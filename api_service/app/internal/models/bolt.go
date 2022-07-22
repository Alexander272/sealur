package models

type BoltDTO struct {
	Title    string  `json:"title"`
	Diameter int32   `json:"diameter"`
	Area     float64 `json:"area"`
}
