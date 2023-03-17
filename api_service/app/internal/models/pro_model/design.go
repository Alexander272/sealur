package pro_model

type Jumper struct {
	HasJumper bool   `json:"hasJumper"`
	Code      string `json:"code"`
	Width     string `json:"width"`
}

type Mounting struct {
	HasMounting bool   `json:"hasMounting"`
	Code        string `json:"code"`
}
