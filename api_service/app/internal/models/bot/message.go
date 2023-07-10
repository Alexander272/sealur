package bot

type Message struct {
	//Можно id канала определять тут (точнее передавать его вместе с остальными данными)
	Service Service     `json:"service" binding:"required"`
	Data    MessageData `json:"data" binding:"required"`
}

type Service struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type MessageData struct {
	Date    string `json:"date" binding:"required"`
	Error   string `json:"error" binding:"required"`
	IP      string `json:"ip" binding:"required"`
	URL     string `json:"url" binding:"required"`
	User    string `json:"user"`
	Company string `json:"company"`
	Request string `json:"request"`
}
