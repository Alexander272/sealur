package moment_model

import (
	"strconv"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type CalcCap struct {
	Pressure       string         `json:"pressure"`
	AxialForce     string         `json:"axialForce"`
	Temp           string         `json:"temp"`
	IsWork         bool           `json:"isWork"`
	Flanges        string         `json:"flanges"`
	IsEmbedded     bool           `json:"isEmbedded"`
	TypeB          string         `json:"type"`
	Condition      string         `json:"condition"`
	Calculation    string         `json:"calculation"`
	Gasket         GasketFullData `json:"gasket"`
	Bolts          BoltsData      `json:"bolts"`
	Embed          EmbedData      `json:"embed"`
	FlangeData     Flanges        `json:"flangeData"`
	CapData        Cap            `json:"capData"`
	IsUseWasher    bool           `json:"isUseWasher"`
	Washer         WasherData     `json:"washer"`
	IsNeedFormulas bool           `json:"isNeedFormulas"`
}

type Cap struct {
	Type     string       `json:"type"`
	H        string       `json:"h"`
	Radius   string       `json:"redius"`
	Delta    string       `json:"delta"`
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
}

func (c *CalcCap) NewCap() (cap *moment_api.CalcCapRequest, err error) {
	pressure, err := strconv.ParseFloat(c.Pressure, 64)
	if err != nil {
		return nil, err
	}
	axialForce, err := strconv.Atoi(c.AxialForce)
	if err != nil {
		return nil, err
	}
	temp, err := strconv.ParseFloat(c.Temp, 64)
	if err != nil {
		return nil, err
	}

	flanges := moment_api.CalcCapRequest_Flanges_value[c.Flanges]
	typeB := moment_api.CalcCapRequest_Type_value[c.TypeB]
	condition := moment_api.CalcCapRequest_Condition_value[c.Condition]
	calculation := moment_api.CalcCapRequest_Calcutation_value[c.Calculation]

	flangeData, err := c.FlangeData.NewFlangeData()
	if err != nil {
		return nil, err
	}

	capData, err := c.CapData.NewCapData()
	if err != nil {
		return nil, err
	}

	bolts, err := c.Bolts.NewBolts()
	if err != nil {
		return nil, err
	}

	gasket, err := c.Gasket.NewGasket()
	if err != nil {
		return nil, err
	}

	var embed *moment_api.EmbedData
	if c.IsEmbedded {
		embed, err = c.Embed.NewEmbed()
		if err != nil {
			return nil, err
		}
	}

	var washer []*moment_api.WasherData
	if c.IsUseWasher {
		washer, err = c.Washer.NewWasher()
		if err != nil {
			return nil, err
		}
	}

	cap = &moment_api.CalcCapRequest{
		Pressure:       pressure,
		AxialForce:     int32(axialForce),
		Temp:           temp,
		IsWork:         c.IsWork,
		Flanges:        moment_api.CalcCapRequest_Flanges(flanges),
		IsEmbedded:     c.IsEmbedded,
		Type:           moment_api.CalcCapRequest_Type(typeB),
		Condition:      moment_api.CalcCapRequest_Condition(condition),
		Calculation:    moment_api.CalcCapRequest_Calcutation(calculation),
		IsUseWasher:    c.IsUseWasher,
		IsNeedFormulas: c.IsNeedFormulas,
		FlangeData:     flangeData,
		CapData:        capData,
		Bolts:          bolts,
		Gasket:         gasket,
		Washer:         washer,
		Embed:          embed,
	}
	return cap, nil
}

func (c *Cap) NewCapData() (cap *moment_api.CapData, err error) {
	typeC := moment_api.CapData_Type_value[c.Type]

	h, err := strconv.ParseFloat(c.H, 64)
	if err != nil {
		return nil, err
	}
	var radius, delta float64
	if c.Radius != "" {
		radius, err = strconv.ParseFloat(c.Radius, 64)
		if err != nil {
			return nil, err
		}
	}
	if c.Delta != "" {
		delta, err = strconv.ParseFloat(c.Delta, 64)
		if err != nil {
			return nil, err
		}
	}

	var mat *moment_api.MaterialData
	if c.MarkId == "another" {
		mat, err = c.Material.NewMaterial()
		if err != nil {
			return nil, err
		}
	}

	cap = &moment_api.CapData{
		Type:     moment_api.CapData_Type(typeC),
		H:        h,
		Radius:   radius,
		Delta:    delta,
		MarkId:   c.MarkId,
		Material: mat,
	}

	return cap, nil
}
