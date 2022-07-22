package models

type FlangeDTO struct {
	Title string `json:"title" binding:"required"`
	Short string `json:"short" binding:"required"`
}

type TypeFlangeDTO struct {
	Title string `json:"title" binding:"required"`
}

type FlangeSizeDTO struct {
	StandId string  `json:"standId,omitempty"`
	Pn      float64 `json:"pn,omitempty"`
	D       float64 `json:"d,omitempty"`
	D6      float64 `json:"d6,omitempty"`
	DOut    float64 `json:"dOut,omitempty"`
	H       float64 `json:"h,omitempty"`
	S0      float64 `json:"s0,omitempty"`
	S1      float64 `json:"s1,omitempty"`
	Length  float64 `json:"length,omitempty"`
	Count   int32   `json:"count,omitempty"`
	BoltId  string  `json:"botId"`
}
