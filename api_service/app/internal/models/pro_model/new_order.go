package pro_model

import "github.com/Alexander272/sealur_proto/api/pro/models/position_model"

type Order struct {
	Id        string `json:"id"`
	Count     string `json:"count"`
	UserId    string
	Positions []Position `json:"positions"`
}

type Position struct {
	OrderId  string   `json:"orderId"`
	Count    int64    `json:"count"`
	Title    string   `json:"title"`
	Amount   string   `json:"amount"`
	Info     string   `json:"info"`
	Type     string   `json:"type"`
	SnpData  SnpData  `json:"snpData"`
	PutgData PutgData `json:"putgData"`
	RingData RingData `json:"ringData"`
	KitData  KitData  `json:"kitData"`
}

func (p *Position) Parse() *position_model.FullPosition {
	positionType := position_model.PositionType_value[p.Type]

	return &position_model.FullPosition{
		// Id: p.,
		OrderId:  p.OrderId,
		Count:    p.Count,
		Title:    p.Title,
		Amount:   p.Amount,
		Info:     p.Info,
		Type:     position_model.PositionType(positionType),
		SnpData:  p.SnpData.Parse(),
		PutgData: p.PutgData.Parse(),
		RingData: p.RingData.Parse(),
		KitData:  p.KitData.Parse(),
	}
}
