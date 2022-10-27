package ex_rect_model

import (
	"strconv"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/express_rectangle_model"
)

type Calc struct {
	Pressure       string         `json:"pressure"`
	TestPressure   string         `json:"testPressure"`
	TypeBolt       string         `json:"type"`
	Condition      string         `json:"condition"`
	Bolts          BoltsData      `json:"bolts"`
	Gasket         GasketFullData `json:"gasket"`
	IsNeedFormulas bool           `json:"isNeedFormulas"`
}

type GasketFullData struct {
	GasketId  string `json:"gasketId"`
	EnvId     string `json:"envId"`
	Thickness string `json:"thickness"`
	// bp - Ширина прокладки
	Width string `json:"width"`
	// L2 - Размер прокладки в продольном направлении
	SizeLong string `json:"sizeLong"`
	// B2 - Размер прокладки в поперечном направление
	SizeTrans string     `json:"sizeTrans"`
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

func (f *Calc) Parse() (ex *calc_api.ExpressRectangleRequest, err error) {
	pressure, err := strconv.ParseFloat(f.Pressure, 64)
	if err != nil {
		return nil, err
	}
	var testPressure float64
	if f.TestPressure != "" {
		testPressure, err = strconv.ParseFloat(f.TestPressure, 64)
		if err != nil {
			return nil, err
		}
	}

	typeBolt := calc_api.ExpressRectangleRequest_TypeBolt_value[f.TypeBolt]
	condition := calc_api.ExpressRectangleRequest_Condition_value[f.Condition]

	bolts, err := f.Bolts.Parse()
	if err != nil {
		return nil, err
	}

	gasket, err := f.Gasket.Parse()
	if err != nil {
		return nil, err
	}

	ex = &calc_api.ExpressRectangleRequest{
		Pressure:       pressure,
		TestPressure:   testPressure,
		TypeBolt:       calc_api.ExpressRectangleRequest_TypeBolt(typeBolt),
		Condition:      calc_api.ExpressRectangleRequest_Condition(condition),
		IsNeedFormulas: f.IsNeedFormulas,
		Bolts:          bolts,
		Gasket:         gasket,
	}
	return ex, nil
}

func (m *MaterialData) Parse() (mat *express_rectangle_model.MaterialData, err error) {
	eAt20, err := strconv.ParseFloat(m.EpsilonAt20, 64)
	if err != nil {
		return nil, err
	}
	sAt20, err := strconv.ParseFloat(m.SigmaAt20, 64)
	if err != nil {
		return nil, err
	}

	mat = &express_rectangle_model.MaterialData{
		Title:       m.Title,
		EpsilonAt20: eAt20,
		SigmaAt20:   sAt20,
	}

	return mat, nil
}

func (b *BoltsData) Parse() (bolts *express_rectangle_model.BoltData, err error) {
	bolts = &express_rectangle_model.BoltData{
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

func (g *GasketFullData) Parse() (gasket *express_rectangle_model.GasketData, err error) {
	thickness, err := strconv.ParseFloat(g.Thickness, 64)
	if err != nil {
		return nil, err
	}
	width, err := strconv.ParseFloat(g.Width, 64)
	if err != nil {
		return nil, err
	}
	sizeTrans, err := strconv.ParseFloat(g.SizeTrans, 64)
	if err != nil {
		return nil, err
	}
	sizeLong, err := strconv.ParseFloat(g.SizeLong, 64)
	if err != nil {
		return nil, err
	}
	gasket = &express_rectangle_model.GasketData{
		GasketId:  g.GasketId,
		EnvId:     g.EnvId,
		Thickness: thickness,
		Width:     width,
		SizeTrans: sizeTrans,
		SizeLong:  sizeLong,
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

func (g *GasketData) Parse() (data *express_rectangle_model.GasketData_Data, err error) {
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

	typeG := express_rectangle_model.GasketData_Type_value[g.TypeG]

	data = &express_rectangle_model.GasketData_Data{
		Title:           g.Title,
		Type:            express_rectangle_model.GasketData_Type(typeG),
		Qo:              q,
		M:               m,
		Compression:     compression,
		Epsilon:         e,
		PermissiblePres: permissiblePres,
	}

	return data, nil
}
