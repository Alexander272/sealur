package pro_model

type StandDTO struct {
	Title string `json:"title" binding:"required"`
}
