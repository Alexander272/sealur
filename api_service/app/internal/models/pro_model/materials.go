package pro_model

type MaterialsDTO struct {
	Title   string `json:"title" binding:"required"`
	TypeMat string `json:"typeMat" binding:"required"`
}

type BoltMaterialsDTO struct {
	Title    string `json:"title" binding:"required"`
	FlangeId string `json:"flangeId" binding:"required"`
}

type Material struct {
	Id       string `json:"id"`
	Code     string `json:"code"`
	Title    string `json:"title"`
	ShortEn  string `json:"shortEn"`
	ShortRus string `json:"shortRus"`
}
