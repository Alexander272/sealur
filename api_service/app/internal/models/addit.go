package models

type CreateAdditDTO struct {
	Materials   string `json:"materials" binding:"required"`
	Mod         string `json:"mod" binding:"required"`
	Temperature string `json:"temperature" binding:"required"`
	Mounting    string `json:"mounting" binding:"required"`
	Graphite    string `json:"graphite" binding:"required"`
	Fillers     string `json:"fillers" binding:"required"`
}

type UpdateMatDTO struct {
	Materials string `json:"materials" binding:"required"`
	TypeCh    string `json:"type" binding:"required"`
	Change    string `json:"change"`
}

type UpdateModDTO struct {
	Mod    string `json:"mod" binding:"required"`
	TypeCh string `json:"type" binding:"required"`
	Change string `json:"change"`
}

type UpdateTempDTO struct {
	Temperature string `json:"temperature" binding:"required"`
	TypeCh      string `json:"type" binding:"required"`
	Change      string `json:"change"`
}

type UpdateMounDTO struct {
	Mounting string `json:"mounting" binding:"required"`
	TypeCh   string `json:"type" binding:"required"`
	Change   string `json:"change"`
}

type UpdateGrapDTO struct {
	Graphite string `json:"graphite" binding:"required"`
	TypeCh   string `json:"type" binding:"required"`
	Change   string `json:"change"`
}

type UpdateFillersDTO struct {
	Fillers string `json:"fillers" binding:"required"`
	TypeCh  string `json:"type" binding:"required"`
	Change  string `json:"change"`
}
