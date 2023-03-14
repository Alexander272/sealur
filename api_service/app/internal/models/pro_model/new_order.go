package pro_model

type Order struct {
	Id        string `json:"id"`
	Count     string `json:"count"`
	UserId    string
	Positions []Position `json:"positions"`
}

type Position struct {
	Count   int64   `json:"count"`
	Title   string  `json:"title"`
	Amount  string  `json:"amount"`
	Type    string  `json:"type"`
	SnpData SnpData `json:"snpData"`
}

// func (p *Position) Parse() *position_model.FullPosition {
// 	return &position_model.FullPosition{
// 		// Id: p.,
// 		Count: p.Count,
// 		Title: p.Title,
// 		Amount: p.Amount,
// 		Type: p.Type,
// 	}
// }
