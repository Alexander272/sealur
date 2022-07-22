package models

type MaterialsDTO struct {
	Title   string `json:"title" binding:"required"`
	TypeMat string `json:"typeMat" binding:"required"`
}

type BoltMaterialsDTO struct {
	Title    string `json:"title" binding:"required"`
	FlangeId string `json:"flangeId" binding:"required"`
}

type MomentMaterialsDTO struct {
	Title string `json:"title" binding:"required"`
}

type Alpha struct {
	Temperature float64 `json:"temperature"`
	Alpha       float64 `json:"alpha"`
}

type CreateAlphaDTO struct {
	MarkId string  `json:"markId" binding:"required"`
	Alpha  []Alpha `json:"alpha"`
}

type UpdateAlphaDTO struct {
	MarkId      string  `json:"markId"`
	Temperature float64 `json:"temperature"`
	Alpha       float64 `json:"alpha"`
}

type Elasticity struct {
	Temperature float64 `json:"temperature"`
	Elasticity  float64 `json:"elasticity"`
}

type CreateElasticityDTO struct {
	MarkId     string       `json:"markId" binding:"required"`
	Elasticity []Elasticity `json:"elasticity"`
}

type UpdateElasticityDTO struct {
	MarkId      string  `json:"markId"`
	Temperature float64 `json:"temperature"`
	Elasticity  float64 `json:"elasticity"`
}

type Voltage struct {
	Temperature float64 `json:"temperature"`
	Voltage     float64 `json:"voltage"`
}

type CreateVoltageDTO struct {
	MarkId  string    `json:"markId" binding:"required"`
	Voltage []Voltage `json:"voltage"`
}

type UpdateVoltageDTO struct {
	MarkId      string  `json:"markId"`
	Temperature float64 `json:"temperature"`
	Voltage     float64 `json:"voltage"`
}
