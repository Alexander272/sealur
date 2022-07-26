package pro_model

type FlangeDTO struct {
	Title string `json:"title" binding:"required"`
	Short string `json:"short" binding:"required"`
}

type TypeFlangeDTO struct {
	Title string `json:"title" binding:"required"`
}
