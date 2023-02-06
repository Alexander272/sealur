package float_model

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/float_model"
)

type Calc struct {
	Pressure       string         `json:"pressure"`
	IsWork         bool           `json:"isWork"`
	HasThorn       bool           `json:"hasThorn"`
	TypeB          string         `json:"type"`
	Condition      string         `json:"condition"`
	Gasket         GasketFullData `json:"gasket"`
	Bolts          BoltsData      `json:"bolts"`
	FlangeData     Flange         `json:"flangeData"`
	CapData        Cap            `json:"capData"`
	IsNeedFormulas bool           `json:"isNeedFormulas"`
}

type Cap struct {
	H        string       `json:"h"`
	Radius   string       `json:"radius"`
	S        string       `json:"s"`
	T        string       `json:"t"`
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
}

type GasketFullData struct {
	GasketId  string     `json:"gasketId"`
	EnvId     string     `json:"envId"`
	Thickness string     `json:"thickness"`
	DOut      string     `json:"d_out"`
	DIn       string     `json:"d_in"`
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
	BoltId   string       `json:"boltId"`
	Distance string       `json:"distance"`
	Temp     string       `json:"temp"`
	Count    string       `json:"count"`
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
	Diameter string       `json:"diameter"`
	Area     string       `json:"area"`
	Title    string       `json:"title"`
}

type MaterialData struct {
	Title       string `json:"title"`
	EpsilonAt20 string `json:"epsilonAt20"`
	Epsilon     string `json:"epsilon"`
	SigmaAt20   string `json:"sigmaAt20"`
	Sigma       string `json:"sigma"`
}

type Flange struct {
	DOut     string       `json:"dOut"`
	D        string       `json:"d"`
	D6       string       `json:"d6"`
	T        string       `json:"t"`
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
	Width    string       `json:"width"`
	DIn      string       `json:"dIn"`
}

func (f *Calc) New() (float *calc_api.FloatRequest, err error) {
	pressure, err := strconv.ParseFloat(f.Pressure, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pressure. error: %w", err)
	}

	condition := calc_api.FloatRequest_Condition_value[f.Condition]
	typeB := calc_api.FlangeRequest_Type_value[f.TypeB]

	flangeData, err := f.FlangeData.NewFlangeData()
	if err != nil {
		return nil, fmt.Errorf("failed to parse flange data. error: %w", err)
	}

	capData, err := f.CapData.NewCapData()
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap data. error: %w", err)
	}

	bolts, err := f.Bolts.NewBolts()
	if err != nil {
		return nil, fmt.Errorf("failed to parse bolts. error: %w", err)
	}

	gasket, err := f.Gasket.NewGasket()
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket. error: %w", err)
	}

	float = &calc_api.FloatRequest{
		Pressure:       pressure,
		HasThorn:       f.HasThorn,
		IsWork:         f.IsWork,
		Type:           calc_api.FloatRequest_Type(typeB),
		Condition:      calc_api.FloatRequest_Condition(condition),
		IsNeedFormulas: f.IsNeedFormulas,
		FlangeData:     flangeData,
		CapData:        capData,
		Bolts:          bolts,
		Gasket:         gasket,
	}
	return float, nil
}

func (f *Flange) NewFlangeData() (flange *float_model.FlangeData, err error) {
	var mat *float_model.MaterialData
	if f.MarkId == "another" {
		mat, err = f.Material.NewMaterial()
		if err != nil {
			return nil, fmt.Errorf("failed to parse flange materials. error: %w", err)
		}
	}

	dOut, err := strconv.ParseFloat(f.DOut, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flange dOut. error: %w", err)
	}
	d, err := strconv.ParseFloat(f.D, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flange d. error: %w", err)
	}
	d6, err := strconv.ParseFloat(f.D6, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flange d6. error: %w", err)
	}
	t, err := strconv.ParseFloat(f.T, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flange t. error: %w", err)
	}
	var width, dIn float64
	if f.Width != "" {
		width, err = strconv.ParseFloat(f.Width, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse flange width. error: %w", err)
		}
	}
	if f.DIn != "" {
		dIn, err = strconv.ParseFloat(f.DIn, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse flange dIn. error: %w", err)
		}
	}

	flange = &float_model.FlangeData{
		DOut:     dOut,
		D:        d,
		D6:       d6,
		T:        t,
		MarkId:   f.MarkId,
		Material: mat,
		Width:    width,
		DIn:      dIn,
	}

	return flange, nil
}

func (c *Cap) NewCapData() (cap *float_model.CapData, err error) {
	h, err := strconv.ParseFloat(c.H, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap h. error: %w", err)
	}

	radius, err := strconv.ParseFloat(c.Radius, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap radius. error: %w", err)
	}

	s, err := strconv.ParseFloat(c.S, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap s. error: %w", err)
	}

	t, err := strconv.ParseFloat(c.T, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap t. error: %w", err)
	}

	var mat *float_model.MaterialData
	if c.MarkId == "another" {
		mat, err = c.Material.NewMaterial()
		if err != nil {
			return nil, fmt.Errorf("failed to parse cap materials. error: %w", err)
		}
	}

	cap = &float_model.CapData{
		H:        h,
		Radius:   radius,
		S:        s,
		T:        t,
		MarkId:   c.MarkId,
		Material: mat,
	}

	return cap, nil
}

func (m *MaterialData) NewMaterial() (mat *float_model.MaterialData, err error) {
	eAt20, err := strconv.ParseFloat(m.EpsilonAt20, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse material epsilon at 20. error: %w", err)
	}
	e, err := strconv.ParseFloat(m.Epsilon, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse material epsilon. error: %w", err)
	}

	var sAt20, s float64
	if m.SigmaAt20 != "" {
		sAt20, err = strconv.ParseFloat(m.SigmaAt20, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse material sigma at 20. error: %w", err)
		}
	}

	if m.Sigma != "" {
		s, err = strconv.ParseFloat(m.Sigma, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse material sigma. error: %w", err)
		}
	}

	mat = &float_model.MaterialData{
		Title:       m.Title,
		EpsilonAt20: eAt20,
		Epsilon:     e,
		SigmaAt20:   sAt20,
		Sigma:       s,
	}

	return mat, nil
}

func (b *BoltsData) NewBolts() (bolts *float_model.BoltData, err error) {
	bolts = &float_model.BoltData{
		MarkId: b.MarkId,
		BoltId: b.BoltId,
		Title:  b.Title,
	}

	distance, err := strconv.ParseFloat(b.Distance, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bolt distance. error: %w", err)
	}
	bolts.Distance = distance

	count, err := strconv.Atoi(b.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bolt count. error: %w", err)
	}
	bolts.Count = int32(count)

	temp, err := strconv.ParseFloat(b.Temp, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bolt temp. error: %w", err)
	}
	bolts.Temp = temp

	if b.BoltId == "another" {
		diameter, err := strconv.ParseFloat(b.Diameter, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt diameter. error: %w", err)
		}
		area, err := strconv.ParseFloat(b.Area, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt area. error: %w", err)
		}
		bolts.Diameter = diameter
		bolts.Area = area
	}
	if b.MarkId == "another" {
		bolts.Material, err = b.Material.NewMaterial()
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt materials. error: %w", err)
		}
	}

	return bolts, nil
}

func (g *GasketFullData) NewGasket() (gasket *float_model.GasketData, err error) {
	thickness, err := strconv.ParseFloat(g.Thickness, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket thickness. error: %w", err)
	}
	dOut, err := strconv.ParseFloat(g.DOut, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket dOut. error: %w", err)
	}
	dIn, err := strconv.ParseFloat(g.DIn, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket dIn. error: %w", err)
	}
	gasket = &float_model.GasketData{
		GasketId:  g.GasketId,
		EnvId:     g.EnvId,
		Thickness: thickness,
		DOut:      dOut,
		DIn:       dIn,
	}
	if g.GasketId == "another" {
		data, err := g.Data.NewGasketData()
		if err != nil {
			return nil, fmt.Errorf("failed to parse gasket data. error: %w", err)
		}

		gasket.Data = data
	}

	return gasket, nil
}

func (g *GasketData) NewGasketData() (data *float_model.GasketData_Data, err error) {
	q, err := strconv.ParseFloat(g.Qo, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket q. error: %w", err)
	}
	m, err := strconv.ParseFloat(g.M, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket m. error: %w", err)
	}
	compression, err := strconv.ParseFloat(g.Compression, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket compression. error: %w", err)
	}
	e, err := strconv.ParseFloat(g.Epsilon, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket epsilon. error: %w", err)
	}
	permissiblePres, err := strconv.ParseFloat(g.PermissiblePres, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket permissible pres. error: %w", err)
	}

	typeG := float_model.GasketData_Type_value[g.TypeG]

	data = &float_model.GasketData_Data{
		Title:           g.Title,
		Type:            float_model.GasketData_Type(typeG),
		Qo:              q,
		M:               m,
		Compression:     compression,
		Epsilon:         e,
		PermissiblePres: permissiblePres,
	}

	return data, nil
}
