package models

type MaterialsDTO struct {
	Title   string `json:"title" binding:"required"`
	TypeMat string `json:"typeMat" binding:"required"`
}

type BoltMaterialsDTO struct {
	Title    string `json:"title" binding:"required"`
	FlangeId string `json:"flangeId" binding:"required"`
}
