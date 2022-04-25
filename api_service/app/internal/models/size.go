package models

type GetSizesDTO struct {
	Flange   string `json:"flange" binding:"required"`
	TypeFlId string `json:"typeFlId" binding:"required"`
	TypePr   string `json:"typePr" binding:"required"`
	StandId  string `json:"standId" binding:"required"`
}

type SizesDTO struct {
	Flange   string `json:"flange" binding:"required"`
	TypeFlId string `json:"typeFlId" binding:"required"`
	Dn       string `json:"dn" binding:"required"`
	Pn       string `json:"pn" binding:"required"`
	TypePr   string `json:"typePr" binding:"required"`
	StandId  string `json:"standId"`
	D4       string `json:"d4"`
	D3       string `json:"d3" binding:"required"`
	D2       string `json:"d2" binding:"required"`
	D1       string `json:"d1"`
	H        string `json:"h" binding:"required"`
	S2       string `json:"s2"`
	S3       string `json:"s3"`
	Number   int32  `json:"number"`
}
