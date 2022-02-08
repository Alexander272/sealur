package models

type CreateAdditDTO struct {
	Materials   string `json:"materials" binding:"required"`
	Mod         string `json:"mod" binding:"required"`
	Temperature string `json:"temperature" binding:"required"`
	Mounting    string `json:"mounting" binding:"required"`
	Graphite    string `json:"graphite" binding:"required"`
	TypeFl      string `json:"typeFl" binding:"required"`
}

type UpdateMatDTO struct {
	Materials string `json:"materials" binding:"required"`
}

type UpdateModDTO struct {
	Mod string `json:"mod" binding:"mod"`
}

type UpdateTempDTO struct {
	Temperature string `json:"temperature" binding:"required"`
}

type UpdateMounDTO struct {
	Mounting string `json:"mounting" binding:"required"`
}

type UpdateGrapDTO struct {
	Graphite string `json:"graphite" binding:"required"`
}

type UpdateFlDTO struct {
	TypeFl string `json:"typeFl" binding:"required"`
}
