package models

type SNPDTO struct {
	StandId   string `json:"standId" binding:"required"`
	FlangeId  string `json:"flangeId" binding:"required"`
	TypeFlId  string `json:"typeFlId" binding:"required"`
	TypePr    string `json:"typePr" binding:"required"`
	Fillers   string `json:"fillers" binding:"required"`
	Materials string `json:"materials" binding:"required"`
	DefMat    string `json:"defMat" binding:"required"`
	Mounting  string `json:"mounting" binding:"required"`
	Graphite  string `json:"graphite" binding:"required"`
}
