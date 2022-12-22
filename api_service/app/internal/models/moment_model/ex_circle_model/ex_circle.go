package ex_circle_model

import (
	"strconv"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_circle_model"
)

type Calc struct {
	Pressure       string         `json:"pressure"`
	TypeBolt       string         `json:"type"`
	Condition      string         `json:"condition"`
	Bolts          BoltsData      `json:"bolts"`
	Gasket         GasketFullData `json:"gasket"`
	IsNeedFormulas bool           `json:"isNeedFormulas"`
}

type GasketFullData struct {
	GasketId  string     `json:"gasketId"`
	EnvId     string     `json:"envId"`
	Thickness string     `json:"thickness"`
	DOut      string     `json:"dOut"`
	DIn       string     `json:"dIn"`
	Data      GasketData `json:"data"`
}

type GasketData struct {
	Title           string `json:"title"`
	TypeG           string `json:"type"`
	Qo              string `json:"qo"`
	M               string `json:"m"`
	Compression     string `json:"compression"`
	Epsilon         string `json:"epsilon"`
	PermissiblePres string `json:"permissiblePres"`
}

type BoltsData struct {
	Count    string       `json:"count"`
	BoltId   string       `json:"boltId"`
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
	Diameter string       `json:"diameter"`
	Area     string       `json:"area"`
	Title    string       `json:"title"`
}

type MaterialData struct {
	Title       string `json:"title"`
	EpsilonAt20 string `json:"epsilonAt20"`
	SigmaAt20   string `json:"sigmaAt20"`
}

func (f *Calc) Parse() (ex *calc_api.ExpressCircleRequest, err error) {
	pressure, err := strconv.ParseFloat(f.Pressure, 64)
	if err != nil {
		return nil, err
	}

	typeBolt := calc_api.ExpressCircleRequest_TypeBolt_value[f.TypeBolt]
	condition := calc_api.ExpressCircleRequest_Condition_value[f.Condition]

	bolts, err := f.Bolts.Parse()
	if err != nil {
		return nil, err
	}

	gasket, err := f.Gasket.Parse()
	if err != nil {
		return nil, err
	}

	ex = &calc_api.ExpressCircleRequest{
		Pressure:       pressure,
		TypeBolt:       calc_api.ExpressCircleRequest_TypeBolt(typeBolt),
		Condition:      calc_api.ExpressCircleRequest_Condition(condition),
		IsNeedFormulas: f.IsNeedFormulas,
		Bolts:          bolts,
		Gasket:         gasket,
	}
	return ex, nil
}

func (m *MaterialData) Parse() (mat *express_circle_model.MaterialData, err error) {
	eAt20, err := strconv.ParseFloat(m.EpsilonAt20, 64)
	if err != nil {
		return nil, err
	}
	sAt20, err := strconv.ParseFloat(m.SigmaAt20, 64)
	if err != nil {
		return nil, err
	}

	mat = &express_circle_model.MaterialData{
		Title:       m.Title,
		EpsilonAt20: eAt20,
		SigmaAt20:   sAt20,
	}

	return mat, nil
}

func (b *BoltsData) Parse() (bolts *express_circle_model.BoltData, err error) {
	bolts = &express_circle_model.BoltData{
		MarkId: b.MarkId,
		BoltId: b.BoltId,
	}

	count, err := strconv.Atoi(b.Count)
	if err != nil {
		return nil, err
	}
	bolts.Count = int32(count)

	if b.BoltId == "another" {
		diameter, err := strconv.ParseFloat(b.Diameter, 64)
		if err != nil {
			return nil, err
		}
		area, err := strconv.ParseFloat(b.Area, 64)
		if err != nil {
			return nil, err
		}
		bolts.Diameter = diameter
		bolts.Area = area
	}
	if b.MarkId == "another" {
		bolts.Material, err = b.Material.Parse()
		if err != nil {
			return nil, err
		}
	}

	return bolts, nil
}

func (g *GasketFullData) Parse() (gasket *express_circle_model.GasketData, err error) {
	thickness, err := strconv.ParseFloat(g.Thickness, 64)
	if err != nil {
		return nil, err
	}
	dOut, err := strconv.ParseFloat(g.DOut, 64)
	if err != nil {
		return nil, err
	}
	dIn, err := strconv.ParseFloat(g.DIn, 64)
	if err != nil {
		return nil, err
	}
	gasket = &express_circle_model.GasketData{
		GasketId:  g.GasketId,
		EnvId:     g.EnvId,
		Thickness: thickness,
		DOut:      dOut,
		DIn:       dIn,
	}
	if g.GasketId == "another" {
		data, err := g.Data.Parse()
		if err != nil {
			return nil, err
		}

		gasket.Data = data
	}

	return gasket, nil
}

func (g *GasketData) Parse() (data *express_circle_model.GasketData_Data, err error) {
	q, err := strconv.ParseFloat(g.Qo, 64)
	if err != nil {
		return nil, err
	}
	m, err := strconv.ParseFloat(g.M, 64)
	if err != nil {
		return nil, err
	}
	compression, err := strconv.ParseFloat(g.Compression, 64)
	if err != nil {
		return nil, err
	}
	e, err := strconv.ParseFloat(g.Epsilon, 64)
	if err != nil {
		return nil, err
	}
	permissiblePres, err := strconv.ParseFloat(g.PermissiblePres, 64)
	if err != nil {
		return nil, err
	}

	typeG := express_circle_model.GasketData_Type_value[g.TypeG]

	data = &express_circle_model.GasketData_Data{
		Title:           g.Title,
		Type:            express_circle_model.GasketData_Type(typeG),
		Qo:              q,
		M:               m,
		Compression:     compression,
		Epsilon:         e,
		PermissiblePres: permissiblePres,
	}

	return data, nil
}
