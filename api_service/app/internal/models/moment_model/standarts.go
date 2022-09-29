package moment_model

type StandartDTO struct {
	Title          string   `json:"title"`
	TypeId         string   `json:"typeId"`
	TitleDn        string   `json:"titleDn"`
	TitlePn        string   `json:"titlePn"`
	IsNeedRow      bool     `json:"isNeedRow"`
	Rows           []string `json:"rows"`
	IsInch         bool     `json:"isInch"`
	HasDesignation bool     `json:"hasDesignation"`
}

type TypeFlangeDTO struct {
	Title string `json:"title"`
	Label string `json:"label"`
}
