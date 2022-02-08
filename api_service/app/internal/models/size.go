package models

type GetSizesDTO struct {
	Flange  string `json:"flange" binding:"required"`
	TypeFl  string `json:"typeFl" binding:"required"`
	TypePr  string `json:"typePr" binding:"required"`
	StandId string `json:"standId" binding:"required"`
}

type DeleteSizeDTO struct {
	Flange string `json:"flange" binding:"required"`
	TypeFl string `json:"typeFl" binding:"required"`
}

type SizesDTO struct {
	Flange  string `json:"flange" binding:"required"`
	TypeFl  string `json:"typeFl" binding:"required"`
	Dn      int32  `json:"dn" binding:"required"`
	Pn      int32  `json:"pn" binding:"required"`
	TypePr  string `json:"typePr" binding:"required"`
	StandId string `json:"standId" binding:"required"`
	D4      int32  `json:"d4" binding:"required"`
	D3      int32  `json:"d3" binding:"required"`
	D2      int32  `json:"d2" binding:"required"`
	D1      int32  `json:"d1" binding:"required"`
	H       int32  `json:"h" binding:"required"`
}
